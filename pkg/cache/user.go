package cache

import (
	"context"
	"encoding/json"
	"lotteryBackend/global"

	"github.com/redis/go-redis/v9"
)

// UserInfoCache 用户信息缓存结构
type UserInfoCache struct {
	ID         uint   `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RealName   string `json:"realName"`
	Department string `json:"department"`
}

// SetUserInfoCache 缓存用户信息
func SetUserInfoCache(userId uint, info *UserInfoCache) error {
	ctx := context.Background()
	key := GetUserInfoKey(userId)

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	return global.GVA_REDIS.Set(ctx, key, data, UserInfoExpire).Err()
}

// GetUserInfoCache 获取用户信息缓存
func GetUserInfoCache(userId uint) (*UserInfoCache, error) {
	ctx := context.Background()
	key := GetUserInfoKey(userId)

	data, err := global.GVA_REDIS.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var cache UserInfoCache
	if err := json.Unmarshal([]byte(data), &cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

// BatchGetUserInfoCache 批量获取用户信息缓存
// 返回：已缓存的用户Map, 未命中的用户ID列表, 错误
func BatchGetUserInfoCache(userIds []uint) (map[uint]*UserInfoCache, []uint, error) {
	ctx := context.Background()
	result := make(map[uint]*UserInfoCache)
	var missedIds []uint

	if len(userIds) == 0 {
		return result, missedIds, nil
	}

	// 使用 Pipeline 批量获取
	pipe := global.GVA_REDIS.Pipeline()
	cmds := make(map[uint]*redis.StringCmd)

	for _, userId := range userIds {
		key := GetUserInfoKey(userId)
		cmds[userId] = pipe.Get(ctx, key)
	}

	pipe.Exec(ctx)

	for userId, cmd := range cmds {
		data, err := cmd.Result()
		if err == redis.Nil {
			missedIds = append(missedIds, userId)
			continue
		}
		if err != nil {
			missedIds = append(missedIds, userId)
			continue
		}

		var cache UserInfoCache
		if err := json.Unmarshal([]byte(data), &cache); err != nil {
			missedIds = append(missedIds, userId)
			continue
		}
		result[userId] = &cache
	}

	return result, missedIds, nil
}

// BatchSetUserInfoCache 批量设置用户信息缓存
func BatchSetUserInfoCache(users map[uint]*UserInfoCache) error {
	if len(users) == 0 {
		return nil
	}

	ctx := context.Background()
	pipe := global.GVA_REDIS.Pipeline()

	for userId, info := range users {
		key := GetUserInfoKey(userId)
		data, _ := json.Marshal(info)
		pipe.Set(ctx, key, data, UserInfoExpire)
	}

	_, err := pipe.Exec(ctx)
	return err
}
