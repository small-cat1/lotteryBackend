package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"lotteryBackend/global"
	"time"

	"github.com/redis/go-redis/v9"
)

// RoundInfoCache 场次信息缓存
type RoundInfoCache struct {
	ID          uint `json:"id"`
	ActivityId  uint `json:"activityId"`
	WinnerCount int  `json:"winnerCount"`
	Duration    int  `json:"duration"`
}

// ==================== 当前场次 ====================

// SetCurrentRound 设置当前进行中的场次
func SetCurrentRound(activityId, roundId uint, duration int) error {
	ctx := context.Background()
	key := GetCurrentRoundKey(activityId)
	// 过期时间 = 游戏时长 + 1小时（预留结算时间）
	expire := time.Duration(duration)*time.Second + time.Hour
	return global.GVA_REDIS.Set(ctx, key, roundId, expire).Err()
}

// SetCurrentRoundNX 原子设置当前场次（防止并发开启多个游戏）
// 返回：是否设置成功
func SetCurrentRoundNX(activityId, roundId uint, duration int) (bool, error) {
	ctx := context.Background()
	key := GetCurrentRoundKey(activityId)
	expire := time.Duration(duration)*time.Second + time.Hour
	return global.GVA_REDIS.SetNX(ctx, key, roundId, expire).Result()
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
	ctx := context.Background()
	key := GetCurrentRoundKey(activityId)

	ttl, err := global.GVA_REDIS.TTL(ctx, key).Result()
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

// ==================== 场次状态 ====================

// SetRoundStatus 设置场次状态
func SetRoundStatus(roundId uint, status int) error {
	ctx := context.Background()
	key := GetRoundStatusKey(roundId)
	return global.GVA_REDIS.Set(ctx, key, status, RoundInfoExpire).Err()
}

// GetRoundStatus 获取场次状态
// 返回：状态值，-1 表示不存在
func GetRoundStatus(roundId uint) (int, error) {
	ctx := context.Background()
	key := GetRoundStatusKey(roundId)

	val, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}

	var status int
	fmt.Sscanf(val, "%d", &status)
	return status, nil
}

// ==================== 场次信息 ====================

// SetRoundInfo 缓存场次信息
func SetRoundInfo(roundId uint, info *RoundInfoCache) error {
	ctx := context.Background()
	key := GetRoundInfoKey(roundId)

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return global.GVA_REDIS.Set(ctx, key, data, RoundInfoExpire).Err()
}

// GetRoundInfo 获取场次信息
func GetRoundInfo(roundId uint) (*RoundInfoCache, error) {
	ctx := context.Background()
	key := GetRoundInfoKey(roundId)

	data, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var info RoundInfoCache
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ==================== 清理函数 ====================

// ClearRoundStatus 清除场次状态
func ClearRoundStatus(roundId uint) error {
	ctx := context.Background()
	key := GetRoundStatusKey(roundId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

// ClearRoundInfo 清除场次信息
func ClearRoundInfo(roundId uint) error {
	ctx := context.Background()
	key := GetRoundInfoKey(roundId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

// ClearRoundAll 清除场次所有缓存（状态、信息、积分）
// ⭐ 游戏开始前调用，防止脏数据
func ClearRoundAll(roundId uint) error {
	ctx := context.Background()
	keys := []string{
		GetRoundStatusKey(roundId),
		GetRoundInfoKey(roundId),
		GetRoundScoresKey(roundId),
	}
	return global.GVA_REDIS.Del(ctx, keys...).Err()
}

// InitRound 初始化场次缓存（游戏开始时调用）
// 1. 清除旧数据
// 2. 设置场次状态
// 3. 设置场次信息
func InitRound(roundId uint, status int, info *RoundInfoCache) error {
	// 先清除旧数据
	if err := ClearRoundAll(roundId); err != nil {
		return err
	}

	// 设置状态
	if err := SetRoundStatus(roundId, status); err != nil {
		return err
	}

	// 设置信息
	if info != nil {
		if err := SetRoundInfo(roundId, info); err != nil {
			return err
		}
	}

	return nil
}
