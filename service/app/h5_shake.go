package app

import (
	"errors"
	"lotteryBackend/constants"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
	"lotteryBackend/pkg/cache"
	"time"
)

type H5ShakeService struct{}

// GetCurrentRound 获取当前进行中的场次
func (s *H5ShakeService) GetCurrentRound(activityId, userId uint) (*response.ShakeRoundResp, error) {
	// 1. 先从缓存获取当前场次ID
	roundId, err := cache.GetCurrentRound(activityId)
	if err != nil {
		return nil, err
	}

	// 缓存中没有进行中的场次
	if roundId == 0 {
		return nil, nil
	}

	// 2. 再从缓存获取场次状态（双重确认）
	status, err := cache.GetRoundStatus(roundId)
	if err != nil || status != constants.ShakeRoundStatusOngoing {
		// 状态不是进行中，返回空
		return nil, nil
	}

	// 3. 场次详情可以从数据库获取（或者也缓存起来）
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, nil
	}
	resp := s.toRoundResp(&round)

	// ⭐ 新增：获取当前用户的分数（从 Redis）
	if userId > 0 {
		myScore, _ := cache.GetUserScore(roundId, userId)
		resp.MyScore = int(myScore)
	}

	return resp, nil
}

// GetShakeRanking 获取实时排名
func (s *H5ShakeService) GetShakeRanking(roundId uint, limit int) ([]response.ShakeRankingResp, error) {
	if limit <= 0 {
		limit = 10
	}

	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	var scores []annual.AnnualShakeScore
	err := global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Limit(limit).
		Find(&scores).Error

	if err != nil {
		return nil, err
	}

	userService := H5UserService{}
	result := make([]response.ShakeRankingResp, len(scores))
	for i, s := range scores {
		userBrief := userService.GetUserBrief(s.UserId)
		result[i] = response.ShakeRankingResp{
			UserId:   s.UserId,
			Score:    s.Score,
			Rank:     i + 1,
			IsWinner: (i + 1) <= round.WinnerCount,
		}
		if userBrief != nil {
			result[i].User = *userBrief
		}
	}

	return result, nil
}

// GetMyScore 获取我的成绩
func (s *H5ShakeService) GetMyScore(userId, roundId uint) (*response.MyScoreResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	var score annual.AnnualShakeScore
	result := global.GVA_DB.Where("round_id = ? AND user_id = ?", roundId, userId).First(&score)

	if result.RowsAffected == 0 {
		return &response.MyScoreResp{
			Score:    0,
			Rank:     0,
			IsWinner: false,
		}, nil
	}

	rank := s.calculateRank(roundId, userId)
	isWinner := rank > 0 && rank <= round.WinnerCount

	return &response.MyScoreResp{
		Score:    score.Score,
		Rank:     rank,
		IsWinner: isWinner,
	}, nil
}

// GetRoundResult 获取场次结果
func (s *H5ShakeService) GetRoundResult(userId, roundId uint) (*response.RoundResultResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	// 获取排名
	ranking, _ := s.GetShakeRanking(roundId, 20)

	// 获取我的成绩
	myScore, _ := s.GetMyScore(userId, roundId)

	// 检查是否中奖
	var winInfo *response.WinningResp
	var winner annual.AnnualWinner
	result := global.GVA_DB.Where("round_id = ? AND user_id = ?", roundId, userId).First(&winner)
	if result.RowsAffected > 0 {
		prizeService := H5PrizeService{}
		winInfo, _ = prizeService.GetWinningDetail(winner.ID)
	}

	return &response.RoundResultResp{
		Ranking:  ranking,
		MyRank:   myScore.Rank,
		IsWinner: myScore.IsWinner,
		WinInfo:  winInfo,
	}, nil
}

// calculateRank 计算用户排名
func (s *H5ShakeService) calculateRank(roundId, userId uint) int {
	var scores []annual.AnnualShakeScore
	global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Find(&scores)

	for i, score := range scores {
		if score.UserId == userId {
			return i + 1
		}
	}

	return 0
}

// toRoundResp 转换为响应结构
func (s *H5ShakeService) toRoundResp(round *annual.AnnualShakeRound) *response.ShakeRoundResp {
	resp := &response.ShakeRoundResp{
		ID:          round.ID,
		RoundName:   round.RoundName,
		Duration:    round.Duration,
		WinnerCount: round.WinnerCount,
		Status:      *round.Status,
		StartTime:   round.StartTime,
		EndTime:     round.EndTime,
	}

	// ⭐ 计算 EndTimeMs（游戏进行中时）
	if *round.Status == 1 && round.StartTime != nil {
		endTime := round.StartTime.Add(time.Duration(round.Duration) * time.Second)
		resp.EndTimeMs = endTime.UnixMilli()
	}

	// 获取关联奖品
	if round.PrizeId > 0 {
		var prize annual.AnnualPrize
		if err := global.GVA_DB.First(&prize, round.PrizeId).Error; err == nil {
			resp.Prize = &response.H5PrizeBrief{
				ID:    prize.ID,
				Name:  prize.Name,
				Image: prize.Image,
				Level: *prize.Level,
			}
		}
	}

	return resp
}
