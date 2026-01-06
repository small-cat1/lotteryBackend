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

// IncrUserScore 增加用户积分
func IncrUserScore(roundId uint, userId uint, score float64) (float64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	newScore, err := global.GVA_REDIS.ZIncrBy(ctx, key, score, fmt.Sprintf("%d", userId)).Result()
	if err != nil {
		return 0, err
	}

	// 设置过期时间
	global.GVA_REDIS.Expire(ctx, key, RoundScoresExpire)

	return newScore, nil
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

// GetUserRank 获取用户排名（从0开始）
func GetUserRank(roundId uint, userId uint) (int64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	rank, err := global.GVA_REDIS.ZRevRank(ctx, key, fmt.Sprintf("%d", userId)).Result()
	if err == redis.Nil {
		return -1, nil // 不在排行榜中
	}
	return rank, err
}

// GetRanking 获取排行榜
func GetRanking(roundId uint, limit int64) ([]RankingItem, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)

	results, err := global.GVA_REDIS.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
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

// GetPlayerCount 获取参与人数
func GetPlayerCount(roundId uint) (int64, error) {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)
	return global.GVA_REDIS.ZCard(ctx, key).Result()
}

// ClearRoundScores 清除场次积分
func ClearRoundScores(roundId uint) error {
	ctx := context.Background()
	key := GetRoundScoresKey(roundId)
	return global.GVA_REDIS.Del(ctx, key).Err()
}
