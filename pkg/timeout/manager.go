package timeout

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Manager 超时管理器
// 负责管理多种类型的超时任务和对应的处理器
type Manager struct {
	engine   TimeoutEngine             // 超时引擎
	handlers map[string]TimeoutHandler // 超时处理器映射: type -> handler
	config   *Config                   // 配置
	logger   *zap.Logger               // 日志
	ctx      context.Context           // 上下文
	cancel   context.CancelFunc        // 取消函数
	wg       sync.WaitGroup            // 等待组
	mu       sync.RWMutex              // 读写锁
}

var (
	globalManager *Manager
	managerOnce   sync.Once
)

// NewManager 创建超时管理器
func NewManager(engine TimeoutEngine, config *Config, logger *zap.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	if config == nil {
		config = DefaultConfig()
	}

	return &Manager{
		engine:   engine,
		handlers: make(map[string]TimeoutHandler),
		config:   config,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// InitGlobalManager 初始化全局管理器（单例）
func InitGlobalManager(engine TimeoutEngine, config *Config, logger *zap.Logger) {
	managerOnce.Do(func() {
		globalManager = NewManager(engine, config, logger)
		logger.Info("全局超时管理器初始化成功", zap.String("engine", engine.Name()))
	})
}

// GetGlobalManager 获取全局管理器
func GetGlobalManager() *Manager {
	if globalManager == nil {
		panic("超时管理器未初始化，请先调用 InitGlobalManager")
	}
	return globalManager
}

// RegisterHandler 注册超时处理器
func (m *Manager) RegisterHandler(handler TimeoutHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	handlerType := handler.Type()
	if _, exists := m.handlers[handlerType]; exists {
		m.logger.Warn("超时处理器已存在，将被覆盖", zap.String("type", handlerType))
	}

	m.handlers[handlerType] = handler
	m.logger.Info("注册超时处理器", zap.String("type", handlerType))
}

// UnregisterHandler 注销超时处理器
func (m *Manager) UnregisterHandler(handlerType string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.handlers, handlerType)
	m.logger.Info("注销超时处理器", zap.String("type", handlerType))
}

// GetHandler 获取超时处理器
func (m *Manager) GetHandler(handlerType string) (TimeoutHandler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	handler, exists := m.handlers[handlerType]
	return handler, exists
}

// Add 添加超时任务
func (m *Manager) Add(handlerType, key string, expireAt time.Time) error {
	// 检查处理器是否存在
	if _, exists := m.GetHandler(handlerType); !exists {
		return fmt.Errorf("超时处理器不存在: %s", handlerType)
	}

	// 将类型和key组合
	fullKey := m.buildFullKey(handlerType, key)

	// 添加到引擎
	if err := m.engine.Add(m.ctx, fullKey, expireAt); err != nil {
		m.logger.Error("添加超时任务失败",
			zap.String("type", handlerType),
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	m.logger.Debug("添加超时任务成功",
		zap.String("type", handlerType),
		zap.String("key", key),
		zap.Any("expireAt", expireAt))

	return nil
}

// Remove 移除超时任务
func (m *Manager) Remove(handlerType, key string) error {
	fullKey := m.buildFullKey(handlerType, key)

	if err := m.engine.Remove(m.ctx, fullKey); err != nil {
		m.logger.Error("移除超时任务失败",
			zap.String("type", handlerType),
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	m.logger.Debug("移除超时任务成功",
		zap.String("type", handlerType),
		zap.String("key", key))

	return nil
}

// Start 启动超时扫描
func (m *Manager) Start() {
	m.logger.Info("超时管理器启动中...",
		zap.String("engine", m.engine.Name()),
		zap.Duration("scanInterval", m.config.ScanInterval),
		zap.Int("batchSize", m.config.BatchSize),
		zap.Int("workerNum", m.config.WorkerNum))

	// 启动扫描器
	m.wg.Add(1)
	go m.scanLoop()

	m.logger.Info("超时管理器启动成功 ✅")
}

// Stop 停止超时扫描（优雅关闭）
func (m *Manager) Stop() {
	m.logger.Info("超时管理器停止中...")

	// 取消context
	m.cancel()

	// 等待所有goroutine完成
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	// 最多等待30秒
	select {
	case <-done:
		m.logger.Info("超时管理器已安全停止")
	case <-time.After(30 * time.Second):
		m.logger.Warn("超时管理器停止超时，强制退出")
	}

	// 关闭引擎
	if err := m.engine.Close(); err != nil {
		m.logger.Error("关闭超时引擎失败", zap.Error(err))
	} else {
		m.logger.Info("超时引擎已关闭")
	}
}

// scanLoop 扫描循环
func (m *Manager) scanLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.config.ScanInterval)
	defer ticker.Stop()

	m.logger.Info("超时扫描器已启动",
		zap.Duration("interval", m.config.ScanInterval))

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("超时扫描器已停止")
			return
		case <-ticker.C:
			m.scan()
		}
	}
}

// scan 执行一次扫描
func (m *Manager) scan() {
	// 从引擎获取过期任务
	keys, err := m.engine.Scan(m.ctx, m.config.BatchSize)
	if err != nil {
		m.logger.Error("扫描超时任务失败", zap.Error(err))
		return
	}

	if len(keys) == 0 {
		return
	}

	m.logger.Info("发现超时任务",
		zap.Int("count", len(keys)),
		zap.String("engine", m.engine.Name()))

	// 创建工作池
	taskChan := make(chan string, len(keys))

	// 启动worker
	for i := 0; i < m.config.WorkerNum; i++ {
		m.wg.Add(1)
		go m.worker(i, taskChan)
	}

	// 发送任务
	for _, key := range keys {
		taskChan <- key
	}

	// 关闭任务通道
	close(taskChan)
}

// worker 处理任务
func (m *Manager) worker(id int, taskChan <-chan string) {
	defer m.wg.Done()

	for fullKey := range taskChan {
		// 解析类型和key
		handlerType, key := m.parseFullKey(fullKey)

		// 获取处理器
		handler, exists := m.GetHandler(handlerType)
		if !exists {
			m.logger.Error("处理器不存在，跳过",
				zap.Int("worker", id),
				zap.String("type", handlerType),
				zap.String("key", key))
			continue
		}

		// 处理超时任务（带重试）
		m.handleWithRetry(id, handler, key)
	}
}

// handleWithRetry 带重试的处理
func (m *Manager) handleWithRetry(workerID int, handler TimeoutHandler, key string) {
	var err error

	for i := 0; i <= m.config.RetryTimes; i++ {
		if i > 0 {
			m.logger.Warn("重试处理超时任务",
				zap.Int("worker", workerID),
				zap.String("type", handler.Type()),
				zap.String("key", key),
				zap.Int("retry", i))
			time.Sleep(m.config.RetryInterval)
		}

		// 调用处理器
		err = handler.Handle(m.ctx, key)
		if err == nil {
			m.logger.Info("处理超时任务成功",
				zap.Int("worker", workerID),
				zap.String("type", handler.Type()),
				zap.String("key", key))
			return
		}

		m.logger.Error("处理超时任务失败",
			zap.Int("worker", workerID),
			zap.String("type", handler.Type()),
			zap.String("key", key),
			zap.Int("retry", i),
			zap.Error(err))
	}

	// 所有重试都失败
	m.logger.Error("处理超时任务最终失败，已达最大重试次数",
		zap.Int("worker", workerID),
		zap.String("type", handler.Type()),
		zap.String("key", key),
		zap.Int("retryTimes", m.config.RetryTimes),
		zap.Error(err))
}

// buildFullKey 构建完整key（类型:key）
func (m *Manager) buildFullKey(handlerType, key string) string {
	return fmt.Sprintf("%s:%s", handlerType, key)
}

// parseFullKey 解析完整key
func (m *Manager) parseFullKey(fullKey string) (handlerType, key string) {
	// 查找第一个冒号
	for i := 0; i < len(fullKey); i++ {
		if fullKey[i] == ':' {
			return fullKey[:i], fullKey[i+1:]
		}
	}
	// 如果没有冒号，返回空类型
	return "", fullKey
}

// GetEngine 获取底层引擎（用于扩展）
func (m *Manager) GetEngine() TimeoutEngine {
	return m.engine
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *Config {
	return m.config
}
