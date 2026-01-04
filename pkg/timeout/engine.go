package timeout

import (
	"context"
	"time"
)

// TimeoutTask 超时任务
type TimeoutTask struct {
	Key      string    // 任务唯一标识
	ExpireAt time.Time // 过期时间
	Data     any       // 附加数据（可选）
}

// TimeoutHandler 超时处理器接口
// 业务方实现此接口来处理具体的超时逻辑
type TimeoutHandler interface {
	// Handle 处理超时任务
	// key: 任务唯一标识
	// 返回: 错误信息
	Handle(ctx context.Context, key string) error

	// Type 返回处理器类型（如 "order", "payment"）
	Type() string
}

// Config 超时引擎配置
type Config struct {
	// EngineType 引擎类型: redis, rocketmq, kafka, memory
	EngineType string

	// ScanInterval 扫描间隔
	ScanInterval time.Duration

	// BatchSize 每次扫描的批次大小
	BatchSize int

	// WorkerNum 并发处理的worker数量
	WorkerNum int

	// RetryTimes 失败重试次数
	RetryTimes int

	// RetryInterval 重试间隔
	RetryInterval time.Duration

	// Redis配置（当EngineType为redis时使用）
	RedisConfig *RedisConfig

	// RocketMQ配置（当EngineType为rocketmq时使用）
	// RocketMQConfig *RocketMQConfig

	// Kafka配置（当EngineType为kafka时使用）
	// KafkaConfig *KafkaConfig
}

// RedisConfig Redis引擎配置
type RedisConfig struct {
	// QueueKey Redis队列的key前缀
	QueueKey string

	// LockKeyPrefix 分布式锁的key前缀
	LockKeyPrefix string

	// LockExpire 分布式锁过期时间
	LockExpire time.Duration
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		EngineType:    "redis",
		ScanInterval:  2 * time.Second,
		BatchSize:     100,
		WorkerNum:     5,
		RetryTimes:    3,
		RetryInterval: 1 * time.Second,
		RedisConfig: &RedisConfig{
			QueueKey:      "timeout:queue",
			LockKeyPrefix: "lock:timeout",
			LockExpire:    30 * time.Second,
		},
	}
}

// EngineFactory 引擎工厂接口
type EngineFactory interface {
	// Create 创建引擎实例
	Create(config *Config) (TimeoutEngine, error)
}
