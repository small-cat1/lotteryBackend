package timeout

import (
	"context"
	"time"
)

// TimeoutEngine 超时引擎接口
// 不同的实现可以使用不同的底层存储（Redis、RocketMQ、Kafka、内存队列等）
type TimeoutEngine interface {
	// Add 添加超时任务
	// key: 任务唯一标识（如订单号、支付流水号）
	// expireAt: 过期时间
	Add(ctx context.Context, key string, expireAt time.Time) error

	// Remove 移除超时任务
	// key: 任务唯一标识
	Remove(ctx context.Context, key string) error

	// Scan 扫描过期的任务
	// limit: 每次扫描的最大数量
	// 返回: 过期任务的key列表
	Scan(ctx context.Context, limit int) ([]string, error)

	// Close 关闭引擎，释放资源
	Close() error

	// Name 返回引擎名称
	Name() string
}
