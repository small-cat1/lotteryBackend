package cache

import (
	"context"
	"fmt"
	"lotteryBackend/global"
	"time"
)

func GetEligibleKey(activityId uint) string {
	return fmt.Sprintf("eligible:activity:%d", activityId)
}

// SetUserEligible 设置用户参与资格
func SetUserEligible(activityId, userId uint) error {
	ctx := context.Background()
	key := GetEligibleKey(activityId)
	err := global.GVA_REDIS.SAdd(ctx, key, userId).Err()
	if err != nil {
		return err
	}
	global.GVA_REDIS.Expire(ctx, key, 7*24*time.Hour)
	return nil
}

// RemoveUserEligible 移除用户参与资格
func RemoveUserEligible(activityId, userId uint) error {
	ctx := context.Background()
	key := GetEligibleKey(activityId)
	return global.GVA_REDIS.SRem(ctx, key, userId).Err()
}

// IsUserEligible 检查用户是否有参与资格
func IsUserEligible(activityId, userId uint) bool {
	ctx := context.Background()
	key := GetEligibleKey(activityId)
	result, _ := global.GVA_REDIS.SIsMember(ctx, key, userId).Result()
	return result
}
