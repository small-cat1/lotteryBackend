package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"lotteryBackend/global"
)

// ==================== 签到开关操作 ====================

// GetCheckInSwitch 获取签到开关状态
func GetCheckInSwitch(activityId uint) (bool, error) {
	ctx := context.Background()
	key := GetCheckInSwitchKey(activityId)

	val, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil // 默认关闭
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// SetCheckInSwitch 设置签到开关状态
func SetCheckInSwitch(activityId uint, enabled bool) error {
	ctx := context.Background()
	key := GetCheckInSwitchKey(activityId)

	val := "0"
	if enabled {
		val = "1"
	}
	return global.GVA_REDIS.Set(ctx, key, val, SwitchExpire).Err()
}

// ==================== 弹幕开关操作 ====================

// GetDanmakuSwitch 获取弹幕开关状态
func GetDanmakuSwitch(activityId uint) (bool, error) {
	ctx := context.Background()
	key := GetDanmakuSwitchKey(activityId)

	val, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "1", nil
}

// SetDanmakuSwitch 设置弹幕开关状态
func SetDanmakuSwitch(activityId uint, enabled bool) error {
	ctx := context.Background()
	key := GetDanmakuSwitchKey(activityId)

	val := "0"
	if enabled {
		val = "1"
	}
	return global.GVA_REDIS.Set(ctx, key, val, SwitchExpire).Err()
}
