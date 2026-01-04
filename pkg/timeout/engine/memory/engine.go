package memory

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"lotteryBackend/pkg/timeout"
)

// Engine 内存超时引擎实现（适用于单机、测试环境）
// 使用最小堆实现，性能较高
type Engine struct {
	queue *priorityQueue
	mu    sync.RWMutex
}

// NewEngine 创建内存引擎
func NewEngine() *Engine {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)

	return &Engine{
		queue: &pq,
	}
}

// Add 添加超时任务
func (e *Engine) Add(ctx context.Context, key string, expireAt time.Time) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	item := &item{
		key:      key,
		expireAt: expireAt,
	}

	heap.Push(e.queue, item)
	return nil
}

// Remove 移除超时任务
func (e *Engine) Remove(ctx context.Context, key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 遍历查找并移除
	for i := 0; i < e.queue.Len(); i++ {
		if (*e.queue)[i].key == key {
			heap.Remove(e.queue, i)
			return nil
		}
	}

	return nil
}

// Scan 扫描过期的任务
func (e *Engine) Scan(ctx context.Context, limit int) ([]string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	keys := make([]string, 0, limit)

	// 从堆顶取出过期的任务
	for e.queue.Len() > 0 && len(keys) < limit {
		// 查看堆顶元素
		top := (*e.queue)[0]

		// 如果还没过期，停止扫描
		if top.expireAt.After(now) {
			break
		}

		// 弹出堆顶元素
		item := heap.Pop(e.queue).(*item)
		keys = append(keys, item.key)
	}

	return keys, nil
}

// Close 关闭引擎
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 清空队列
	*e.queue = make(priorityQueue, 0)
	return nil
}

// Name 返回引擎名称
func (e *Engine) Name() string {
	return "memory"
}

// GetQueueLength 获取队列长度（用于监控）
func (e *Engine) GetQueueLength() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.queue.Len()
}

// ============================================
// 优先队列实现（最小堆）
// ============================================

type item struct {
	key      string
	expireAt time.Time
	index    int // 在堆中的索引
}

type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// 过期时间早的优先级高
	return pq[i].expireAt.Before(pq[j].expireAt)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.index = -1 // 标记为已移除
	*pq = old[0 : n-1]
	return item
}

// ============================================
// 引擎工厂
// ============================================

type Factory struct{}

func (f *Factory) Create(config *timeout.Config) (timeout.TimeoutEngine, error) {
	return NewEngine(), nil
}
