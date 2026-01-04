package initialize

import (
	"lotteryBackend/global"
	"lotteryBackend/pkg/timeout"
	"lotteryBackend/pkg/timeout/engine/redisEngine"
	"time"

	"go.uber.org/zap"
)

// Timeout 初始化超时服务
func Timeout() {
	global.GVA_LOG.Info("初始化超时服务...")

	// 1. 检查依赖
	if global.GVA_REDIS == nil {
		panic("Redis未初始化，超时服务无法启动")
		return
	}

	if global.GVA_DB == nil {
		global.GVA_LOG.Error("数据库未初始化，超时服务无法启动")
		return
	}

	// 2. 创建超时引擎（这里使用Redis引擎）
	engine := redisEngine.NewEngine(global.GVA_REDIS, &timeout.RedisConfig{
		QueueKey:      "timeout:queue",
		LockKeyPrefix: "lock:timeout",
		LockExpire:    30 * time.Second,
	})

	// 3. 创建配置
	config := &timeout.Config{
		EngineType:    "redis",
		ScanInterval:  2 * time.Second, // 订单超时扫描间隔
		BatchSize:     100,             // 每次扫描100个
		WorkerNum:     5,               // 5个并发worker
		RetryTimes:    3,               // 失败重试3次
		RetryInterval: 1 * time.Second, // 重试间隔1秒
	}

	// 4. 初始化全局管理器
	timeout.InitGlobalManager(engine, config, global.GVA_LOG)

	// 5. 获取管理器
	manager := timeout.GetGlobalManager()

	// 6. 注册业务处理器

	// 7. 启动超时扫描
	manager.Start()

	global.GVA_LOG.Info("超时服务初始化完成 ✅",
		zap.String("engine", engine.Name()),
		zap.Duration("scanInterval", config.ScanInterval),
		zap.Int("workerNum", config.WorkerNum))
}

// StopTimeout 停止超时服务（优雅关闭时调用）
func StopTimeout() {
	global.GVA_LOG.Info("开始停止超时服务...")

	manager := timeout.GetGlobalManager()
	manager.Stop()

	global.GVA_LOG.Info("超时服务已停止 ✅")
}
