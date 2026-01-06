package ws

import (
	"encoding/json"
	"fmt"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/pkg/cache"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ShakeHandler 摇一摇消息处理器
type ShakeHandler struct {
	lastBroadcastTime map[uint]time.Time
	broadcastMutex    sync.Mutex
}

var shakeHandler = &ShakeHandler{
	lastBroadcastTime: make(map[uint]time.Time),
}

const broadcastInterval = 300 * time.Millisecond

func GetShakeHandler() *ShakeHandler {
	return shakeHandler
}

// HandleShakeScore 处理用户上报分数
func (h *ShakeHandler) HandleShakeScore(client *Client, payload string) {
	var scorePayload ShakeScorePayload
	if err := json.Unmarshal([]byte(payload), &scorePayload); err != nil {
		client.SendMessage(TypeError, ErrorPayload{Code: 400, Message: "分数数据格式错误"})
		return
	}

	userId := client.UserID
	if userId == 0 {
		return
	}

	roundId := scorePayload.RoundId
	score := scorePayload.Score

	// 从 Redis 检查场次状态（比查数据库快）
	status, err := cache.GetRoundStatus(roundId)
	if err != nil || status != 1 {
		return // 游戏未开始或已结束
	}

	// 获取场次信息（用于广播）
	// ⭐ 从缓存获取场次信息
	roundInfo, err := cache.GetRoundInfo(roundId)
	if err != nil {
		// 缓存未命中，回退到数据库
		var round annual.AnnualShakeRound
		if err := global.GVA_DB.Select("id, activity_id, winner_count").First(&round, roundId).Error; err != nil {
			return
		}
		roundInfo = &cache.RoundInfoCache{
			ID:          round.ID,
			ActivityId:  round.ActivityId,
			WinnerCount: round.WinnerCount,
			Duration:    round.Duration,
		}
		// 回填缓存
		go cache.SetRoundInfo(roundId, roundInfo)
	}

	// 更新 Redis 中的分数（使用 ZSet，只保留最高分）
	currentScore, _ := cache.GetUserScore(roundId, userId)
	if float64(score) > currentScore {
		cache.IncrUserScore(roundId, userId, float64(score)-currentScore)

		// 分数变化，广播排名更新
		go h.broadcastRankingUpdate(roundInfo.ActivityId, roundId, roundInfo.WinnerCount)
	}

	global.GVA_LOG.Info("收到摇一摇分数",
		zap.Uint("userId", userId),
		zap.Uint("roundId", roundId),
		zap.Int("score", score))
}

// shouldBroadcast 防抖检查
func (h *ShakeHandler) shouldBroadcast(roundId uint) bool {
	h.broadcastMutex.Lock()
	defer h.broadcastMutex.Unlock()

	now := time.Now()
	if last, ok := h.lastBroadcastTime[roundId]; ok && now.Sub(last) < broadcastInterval {
		return false
	}
	h.lastBroadcastTime[roundId] = now
	return true
}

// broadcastRankingUpdate 广播排名更新到主持人端
func (h *ShakeHandler) broadcastRankingUpdate(activityId, roundId uint, winnerCount int) {
	if !h.shouldBroadcast(roundId) {
		return
	}

	// 1. 从 Redis 获取前10名
	redisRanking, err := cache.GetRanking(roundId, 10)
	if err != nil || len(redisRanking) == 0 {
		return
	}

	// 2. 收集用户ID
	userIds := make([]uint, len(redisRanking))
	for i, r := range redisRanking {
		userIds[i] = r.UserId
	}

	// 3. 从 Redis 批量获取用户信息
	userCache, missedIds, _ := cache.BatchGetUserInfoCache(userIds)

	// 4. 缓存未命中的用户，从 MySQL 查询并回填
	if len(missedIds) > 0 {
		var users []annual.AnnualUser
		global.GVA_DB.Where("id IN ?", missedIds).Find(&users)

		var checkIns []annual.AnnualCheckIn
		global.GVA_DB.Where("activity_id = ? AND user_id IN ?", activityId, missedIds).Find(&checkIns)
		checkInMap := make(map[uint]annual.AnnualCheckIn)
		for _, c := range checkIns {
			checkInMap[c.UserId] = c
		}

		// 回填缓存
		toCache := make(map[uint]*cache.UserInfoCache)
		for _, u := range users {
			checkIn := checkInMap[u.ID]
			cache := &cache.UserInfoCache{
				ID:         u.ID,
				Nickname:   u.Nickname,
				Avatar:     u.Avatar,
				RealName:   checkIn.RealName,
				Department: checkIn.Department,
			}
			userCache[u.ID] = cache
			toCache[u.ID] = cache
		}

		// 异步写入 Redis
		go cache.BatchSetUserInfoCache(toCache)
	}

	// 5. 组装排名数据
	ranking := make([]RankingItem, len(redisRanking))
	for i, r := range redisRanking {
		info := userCache[r.UserId]
		if info == nil {
			info = &cache.UserInfoCache{ID: r.UserId, Nickname: "未知用户"}
		}

		ranking[i] = RankingItem{
			Rank:   r.Rank,
			UserId: r.UserId,
			User: UserBrief{
				ID:         info.ID,
				Nickname:   info.Nickname,
				Avatar:     info.Avatar,
				RealName:   info.RealName,
				Department: info.Department,
			},
			Score:    int(r.Score),
			IsWinner: r.Rank <= winnerCount,
		}
	}

	// 6. 广播
	payload := RankingUpdatePayload{RoundId: roundId, Ranking: ranking}
	GetHub().BroadcastToRoom(fmt.Sprintf("%s:%d", RoomTypeShake, activityId), TypeRankingUpdate, payload)
	GetHub().BroadcastToRoom(fmt.Sprintf("%s:%d", RoomTypeScreen, activityId), TypeRankingUpdate, payload)
}
