package cache

import (
	"context"
	"fmt"
	"lotteryBackend/global"

	"github.com/redis/go-redis/v9"
)

// RankingItem 排行榜项
type RankingItem struct {
	UserId uint    `json:"userId"`
	Score  float64 `json:"score"`
	Rank   int     `json:"rank"`
}

// ==================== 积分操作 ====================

// IncrUserScore 增加用户积分
// 使用 Lua 脚本原子操作：增加分数 + 仅首次设置过期时间
func IncrUserScore(roundId uint, userId uint, score float64) (float64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	script := redis.NewScript(`
		local score = redis.call('ZINCRBY', KEYS[1], ARGV[1], ARGV[2])
		if redis.call('TTL', KEYS[1]) == -1 then
			redis.call('EXPIRE', KEYS[1], ARGV[3])
		end
		return score
	`)

	result, err := script.Run(ctx, global.GVA_REDIS, []string{key},
		score,
		fmt.Sprintf("%d", userId),
		int(RoundScoresExpire.Seconds()),
	).Float64()

	if err != nil {
		return 0, err
	}

	return result, nil
}

// SetUserScore 直接设置用户积分（覆盖）
func SetUserScore(roundId uint, userId uint, score float64) error {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	_, err := global.GVA_REDIS.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: fmt.Sprintf("%d", userId),
	}).Result()

	if err != nil {
		return err
	}

	// 仅首次设置过期时间
	global.GVA_REDIS.ExpireNX(ctx, key, RoundScoresExpire)
	return nil
}

// GetUserScore 获取用户积分
func GetUserScore(roundId uint, userId uint) (float64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	score, err := global.GVA_REDIS.ZScore(ctx, key, fmt.Sprintf("%d", userId)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	return score, err
}

// GetUserRank 获取用户排名（从1开始）
// 返回 -1 表示不在排行榜中
func GetUserRank(roundId uint, userId uint) (int, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	rank, err := global.GVA_REDIS.ZRevRank(ctx, key, fmt.Sprintf("%d", userId)).Result()
	if err == redis.Nil {
		return -1, nil
	}
	if err != nil {
		return -1, err
	}
	return int(rank) + 1, nil // 转为从1开始
}

// ==================== 排行榜 ====================

// GetRanking 获取排行榜
// limit <= 0 表示获取全部
func GetRanking(roundId uint, limit int64) ([]RankingItem, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	var end int64 = -1
	if limit > 0 {
		end = limit - 1
	}

	results, err := global.GVA_REDIS.ZRevRangeWithScores(ctx, key, 0, end).Result()
	if err != nil {
		return nil, err
	}

	ranking := make([]RankingItem, 0, len(results))
	for i, z := range results {
		var userId uint
		fmt.Sscanf(z.Member.(string), "%d", &userId)
		ranking = append(ranking, RankingItem{
			UserId: userId,
			Score:  z.Score,
			Rank:   i + 1,
		})
	}

	return ranking, nil
}

// GetTopN 获取前N名
func GetTopN(roundId uint, n int) ([]RankingItem, error) {
	return GetRanking(roundId, int64(n))
}

// ==================== 统计 ====================

// GetPlayerCount 获取参与人数
func GetPlayerCount(roundId uint) (int64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)
	return global.GVA_REDIS.ZCard(ctx, key).Result()
}

// ==================== 清理 ====================

// ClearRoundScores 清除场次积分
func ClearRoundScores(roundId uint) error {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

// RemoveUserScore 移除用户积分
func RemoveUserScore(roundId uint, userId uint) error {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)
	return global.GVA_REDIS.ZRem(ctx, key, fmt.Sprintf("%d", userId)).Err()
}
