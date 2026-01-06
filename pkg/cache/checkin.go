package cache

import (
	"context"
	"lotteryBackend/global"
)

// AddCheckInUser 添加签到用户
func AddCheckInUser(activityId uint, userId uint) (bool, error) {
	ctx := context.Background()
	key := GetCheckInUsersKey(activityId)

	// SADD 返回添加成功的数量，如果用户已存在返回0
	added, err := global.GVA_REDIS.SAdd(ctx, key, userId).Result()
	if err != nil {
		return false, err
	}

	// 设置过期时间
	global.GVA_REDIS.Expire(ctx, key, CheckInUsersExpire)

	return added > 0, nil
}

// IsUserCheckIn 检查用户是否已签到
func IsUserCheckIn(activityId uint, userId uint) (bool, error) {
	ctx := context.Background()
	key := GetCheckInUsersKey(activityId)
	return global.GVA_REDIS.SIsMember(ctx, key, userId).Result()
}

// GetCheckInCount 获取签到人数
func GetCheckInCount(activityId uint) (int64, error) {
	ctx := context.Background()
	key := GetCheckInUsersKey(activityId)
	return global.GVA_REDIS.SCard(ctx, key).Result()
}

// GetCheckInUserIds 获取所有签到用户ID
func GetCheckInUserIds(activityId uint) ([]string, error) {
	ctx := context.Background()
	key := GetCheckInUsersKey(activityId)
	return global.GVA_REDIS.SMembers(ctx, key).Result()
}
