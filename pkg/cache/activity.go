package cache

import (
	"context"
	"encoding/json"
	"lotteryBackend/global"

	"github.com/redis/go-redis/v9"
)

// ActivityCache 活动缓存结构
type ActivityCache struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Logo           string `json:"logo"`
	Cover          string `json:"cover"`
	Description    string `json:"description"`
	CheckInEnabled int    `json:"checkInEnabled"`
	DanmakuEnabled int    `json:"danmakuEnabled"`
	DanmakuAudit   int    `json:"danmakuAudit"`
	WinnerExclude  int    `json:"winnerExclude"`
	Status         int    `json:"status"`
}

// GetActivityCache 获取活动缓存
func GetActivityCache(activityId uint) (*ActivityCache, error) {
	ctx := context.Background()
	key := GetActivityInfoKey(activityId)

	data, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 缓存不存在
	}
	if err != nil {
		return nil, err
	}

	var cache ActivityCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

// SetActivityCache 设置活动缓存
func SetActivityCache(activity *ActivityCache) error {
	ctx := context.Background()
	key := GetActivityInfoKey(activity.ID)

	data, err := json.Marshal(activity)
	if err != nil {
		return err
	}

	return global.GVA_REDIS.Set(ctx, key, data, ActivityCacheExpire).Err()
}

// DeleteActivityCache 删除活动缓存
func DeleteActivityCache(activityId uint) error {
	ctx := context.Background()
	key := GetActivityInfoKey(activityId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}
