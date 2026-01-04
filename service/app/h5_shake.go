package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
	"sort"
	"time"
)

type H5ShakeService struct{}

// GetCurrentRound 获取当前进行中的场次
func (s *H5ShakeService) GetCurrentRound(activityId uint) (*response.ShakeRoundResp, error) {
	var round annual.AnnualShakeRound

	// 先找进行中的
	result := global.GVA_DB.Where("activity_id = ? AND status = ?", activityId, 1).First(&round)

	if result.RowsAffected == 0 {
		// 再找等待中的
		result = global.GVA_DB.Where("activity_id = ? AND status = ?", activityId, 0).
			Order("sort ASC, id ASC").First(&round)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return s.toRoundResp(&round), nil
}

// GetRoundList 获取场次列表
func (s *H5ShakeService) GetRoundList(activityId uint) ([]response.ShakeRoundResp, error) {
	var rounds []annual.AnnualShakeRound
	err := global.GVA_DB.Where("activity_id = ?", activityId).
		Order("sort ASC, id ASC").
		Find(&rounds).Error

	if err != nil {
		return nil, err
	}

	result := make([]response.ShakeRoundResp, len(rounds))
	for i, r := range rounds {
		result[i] = *s.toRoundResp(&r)
	}

	return result, nil
}

// GetRoundDetail 获取场次详情
func (s *H5ShakeService) GetRoundDetail(roundId uint) (*response.ShakeRoundResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	return s.toRoundResp(&round), nil
}

// JoinGame 加入游戏
func (s *H5ShakeService) JoinGame(userId, roundId uint) error {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return errors.New("场次不存在")
	}

	if *round.Status == 2 {
		return errors.New("本场游戏已结束")
	}

	// 检查是否已经有成绩记录
	var existScore annual.AnnualShakeScore
	result := global.GVA_DB.Where("round_id = ? AND user_id = ?", roundId, userId).First(&existScore)

	if result.RowsAffected == 0 {
		// 创建成绩记录
		score := annual.AnnualShakeScore{
			RoundId: roundId,
			UserId:  userId,
			Score:   0,
		}
		global.GVA_DB.Create(&score)
	}

	return nil
}

// SubmitScore 提交分数
func (s *H5ShakeService) SubmitScore(userId, roundId uint, score int) (*response.MyScoreResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	if *round.Status != 1 {
		return nil, errors.New("游戏未开始或已结束")
	}

	// 更新或创建成绩
	var existScore annual.AnnualShakeScore
	result := global.GVA_DB.Where("round_id = ? AND user_id = ?", roundId, userId).First(&existScore)

	if result.RowsAffected == 0 {
		// 创建
		existScore = annual.AnnualShakeScore{
			RoundId: roundId,
			UserId:  userId,
			Score:   score,
		}
		global.GVA_DB.Create(&existScore)
	} else {
		// 更新（只允许分数增加）
		if score > existScore.Score {
			global.GVA_DB.Model(&existScore).Update("score", score)
			existScore.Score = score
		}
	}

	// 计算排名
	rank := s.calculateRank(roundId, userId)
	isWinner := rank > 0 && rank <= round.WinnerCount

	return &response.MyScoreResp{
		Score:    existScore.Score,
		Rank:     rank,
		IsWinner: isWinner,
	}, nil
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

	// 计算剩余时间
	if *round.Status == 1 && round.StartTime != nil {
		elapsed := int(time.Since(*round.StartTime).Seconds())
		resp.RemainTime = round.Duration - elapsed
		if resp.RemainTime < 0 {
			resp.RemainTime = 0
		}
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

// SettleRound 结算场次（生成中奖记录）
func (s *H5ShakeService) SettleRound(roundId uint) error {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return errors.New("场次不存在")
	}

	// 获取排名
	var scores []annual.AnnualShakeScore
	global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Find(&scores)

	// 按分数排序后取前N名
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	winnerCount := round.WinnerCount
	if winnerCount > len(scores) {
		winnerCount = len(scores)
	}

	// 创建中奖记录
	for i := 0; i < winnerCount; i++ {
		// 更新成绩记录
		global.GVA_DB.Model(&scores[i]).Updates(map[string]interface{}{
			"rank":      i + 1,
			"is_winner": 1,
		})

		// 创建中奖记录
		WinType := 1
		Status := 0
		winner := annual.AnnualWinner{
			ActivityId: round.ActivityId,
			UserId:     scores[i].UserId,
			PrizeId:    round.PrizeId,
			RoundId:    roundId,
			WinType:    &WinType, // 摇一摇
			Status:     &Status,  // 未领取
		}
		global.GVA_DB.Create(&winner)
	}

	// 更新奖品剩余数量
	if round.PrizeId > 0 {
		global.GVA_DB.Model(&annual.AnnualPrize{}).
			Where("id = ?", round.PrizeId).
			Update("remain_count", global.GVA_DB.Raw("remain_count - ?", winnerCount))
	}

	return nil
}
