package cache

import (
	"context"
	"fmt"
	"lotteryBackend/global"
	"time"

	"github.com/redis/go-redis/v9"
)

// SetCurrentRound 设置当前进行中的场次
func SetCurrentRound(activityId, roundId uint, duration int) error {
	key := fmt.Sprintf("annual:activity:%d:current_round", activityId)
	return global.GVA_REDIS.Set(context.Background(), key, roundId, time.Duration(duration)*time.Second).Err()
}

// GetCurrentRound 获取当前进行中的场次
func GetCurrentRound(activityId uint) (uint, error) {
	ctx := context.Background()
	key := GetCurrentRoundKey(activityId)

	val, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	var roundId uint
	fmt.Sscanf(val, "%d", &roundId)
	return roundId, nil
}

// GetCurrentRoundRemaining 获取当前场次剩余时间（秒）
func GetCurrentRoundRemaining(activityId uint) (int, error) {
	key := fmt.Sprintf("annual:activity:%d:current_round", activityId)
	ttl, err := global.GVA_REDIS.TTL(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}
	if ttl < 0 {
		return 0, nil
	}
	return int(ttl.Seconds()), nil
}

// ClearCurrentRound 清除当前场次
func ClearCurrentRound(activityId uint) error {
	ctx := context.Background()
	key := GetCurrentRoundKey(activityId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

// SetRoundStatus 设置场次状态
func SetRoundStatus(roundId uint, status int) error {
	ctx := context.Background()
	key := GetRoundStatusKey(roundId)
	return global.GVA_REDIS.Set(ctx, key, status, RoundScoresExpire).Err()
}

// GetRoundStatus 获取场次状态
func GetRoundStatus(roundId uint) (int, error) {
	ctx := context.Background()
	key := GetRoundStatusKey(roundId)

	val, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return -1, nil // 不存在
	}
	if err != nil {
		return -1, err
	}

	var status int
	fmt.Sscanf(val, "%d", &status)
	return status, nil
}
