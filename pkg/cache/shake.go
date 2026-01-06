package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"lotteryBackend/global"
	"time"
)

// RoundInfoCache 场次信息缓存
type RoundInfoCache struct {
	ID          uint `json:"id"`
	ActivityId  uint `json:"activityId"`
	WinnerCount int  `json:"winnerCount"`
	Duration    int  `json:"duration"`
}

// SetRoundInfo 缓存场次信息
func SetRoundInfo(roundId uint, info *RoundInfoCache) error {
	key := fmt.Sprintf("shake:round:%d:info", roundId)
	data, _ := json.Marshal(info)
	return global.GVA_REDIS.Set(context.Background(), key, data, 2*time.Hour).Err()
}

// GetRoundInfo 获取场次信息
func GetRoundInfo(roundId uint) (*RoundInfoCache, error) {
	key := fmt.Sprintf("shake:round:%d:info", roundId)
	data, err := global.GVA_REDIS.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var info RoundInfoCache
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return &info, nil
}
