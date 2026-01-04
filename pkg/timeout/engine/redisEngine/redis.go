package redisEngine

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"lotteryBackend/pkg/timeout"
	"time"
)

// Engine Redis超时引擎实现
type Engine struct {
	client     redis.UniversalClient
	queueKey   string        // Redis有序集合的key
	lockKey    string        // 分布式锁的key前缀
	lockExpire time.Duration // 锁过期时间
}

// NewEngine 创建Redis引擎
func NewEngine(client redis.UniversalClient, config *timeout.RedisConfig) *Engine {
	if config == nil {
		config = &timeout.RedisConfig{
			QueueKey:      "timeout:queue",
			LockKeyPrefix: "lock:timeout",
			LockExpire:    30 * time.Second,
		}
	}

	return &Engine{
		client:     client,
		queueKey:   config.QueueKey,
		lockKey:    config.LockKeyPrefix,
		lockExpire: config.LockExpire,
	}
}

// Add 添加超时任务
func (e *Engine) Add(ctx context.Context, key string, expireAt time.Time) error {
	score := float64(expireAt.Unix())

	return e.client.ZAdd(ctx, e.queueKey, redis.Z{
		Score:  score,
		Member: key,
	}).Err()
}

// Remove 移除超时任务
func (e *Engine) Remove(ctx context.Context, key string) error {
	return e.client.ZRem(ctx, e.queueKey, key).Err()
}

// Scan 扫描过期的任务
func (e *Engine) Scan(ctx context.Context, limit int) ([]string, error) {
	now := time.Now().Unix()

	// 从有序集合中获取score小于等于当前时间的元素
	keys, err := e.client.ZRangeByScore(ctx, e.queueKey, &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%d", now),
		Count: int64(limit),
	}).Result()

	if err != nil {
		return nil, err
	}

	// 移除已扫描的任务（防止重复处理）
	if len(keys) > 0 {
		// 使用pipeline提高性能
		pipe := e.client.Pipeline()
		for _, key := range keys {
			pipe.ZRem(ctx, e.queueKey, key)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			return nil, fmt.Errorf("移除已扫描任务失败: %w", err)
		}
	}

	return keys, nil
}

// AcquireLock 获取分布式锁（可选，如果业务需要）
func (e *Engine) AcquireLock(ctx context.Context, key string) (bool, error) {
	lockKey := fmt.Sprintf("%s:%s", e.lockKey, key)
	return e.client.SetNX(ctx, lockKey, "1", e.lockExpire).Result()
}

// ReleaseLock 释放分布式锁
func (e *Engine) ReleaseLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("%s:%s", e.lockKey, key)
	return e.client.Del(ctx, lockKey).Err()
}

// Close 关闭引擎
func (e *Engine) Close() error {
	// Redis客户端由外部管理，这里不关闭
	return nil
}

// Name 返回引擎名称
func (e *Engine) Name() string {
	return "redis"
}

// GetQueueLength 获取队列长度（用于监控）
func (e *Engine) GetQueueLength(ctx context.Context) (int64, error) {
	return e.client.ZCard(ctx, e.queueKey).Result()
}

// GetExpiredCount 获取已过期但未处理的任务数量（用于监控）
func (e *Engine) GetExpiredCount(ctx context.Context) (int64, error) {
	now := time.Now().Unix()
	return e.client.ZCount(ctx, e.queueKey, "0", fmt.Sprintf("%d", now)).Result()
}

// Clear 清空队列（用于测试或清理）
func (e *Engine) Clear(ctx context.Context) error {
	return e.client.Del(ctx, e.queueKey).Err()
}
