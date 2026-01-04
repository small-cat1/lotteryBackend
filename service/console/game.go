package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/service/common"
	"time"
)

type GameService struct{}

// 场次状态常量
const (
	RoundStatusPending  = 0 // 未开始
	RoundStatusPlaying  = 1 // 进行中
	RoundStatusFinished = 2 // 已结束
)

// StartGame 开始游戏（验证密码）
func (s *GameService) StartGame(roundId uint, password string) error {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return errors.New("场次不存在")
	}

	// 验证密码
	if round.Password != password {
		return errors.New("密码错误")
	}
	// 检查是否有进行中的游戏
	currentRoundId, err := common.GetCurrentRound(round.ActivityId)
	if err == nil && currentRoundId > 0 {
		return errors.New("当前有进行中的游戏，请先结束后再开始")
	}
	// 更新数据库状态
	now := time.Now()
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":     RoundStatusPlaying,
		"start_time": now,
	})

	// 设置Redis（TTL=游戏时长）
	duration := round.Duration
	if duration <= 0 {
		duration = 60 // 默认1分钟
	}
	common.SetCurrentRound(round.ActivityId, roundId, duration)
	common.SetRoundStatus(roundId, RoundStatusPlaying)

	return nil
}

// GetCurrent 获取当前进行中的游戏
func (s *GameService) GetCurrent(activityId uint) (*response.CurrentGameResp, error) {
	resp := &response.CurrentGameResp{
		Round:   nil,
		Status:  0,
		Ranking: []response.RankingItem{},
	}

	// 获取当前场次ID
	currentRoundId, err := common.GetCurrentRound(activityId)
	if err != nil || currentRoundId == 0 {
		return resp, nil
	}

	// 查询场次详情
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.Preload("Prize").First(&round, currentRoundId).Error; err != nil {
		return resp, nil
	}

	// 构建场次信息
	resp.Round = &response.RoundInfo{
		ID:          round.ID,
		RoundName:   round.RoundName,
		Duration:    round.Duration,
		WinnerCount: round.WinnerCount,
	}

	if round.Prize.ID > 0 {
		resp.Round.Prize = &response.PrizeInfo{
			ID:    round.Prize.ID,
			Name:  round.Prize.Name,
			Image: round.Prize.Image,
			Level: safeInt(round.Prize.Level),
		}
	}

	// 获取剩余时间（直接用TTL）
	remaining, _ := common.GetCurrentRoundRemaining(activityId)
	resp.Remaining = remaining
	resp.Status = RoundStatusPlaying
	// 获取排行榜
	ranking, _ := common.GetRanking(currentRoundId, 20)
	if len(ranking) > 0 {
		userIds := make([]uint, 0, len(ranking))
		for _, item := range ranking {
			userIds = append(userIds, item.UserId)
		}
		var users []annual.AnnualUser
		global.GVA_DB.Where("id IN ?", userIds).Find(&users)
		userMap := make(map[uint]*annual.AnnualUser)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}
		for _, item := range ranking {
			rankItem := response.RankingItem{
				UserId: item.UserId,
				Score:  int(item.Score),
			}
			if user, ok := userMap[item.UserId]; ok {
				rankItem.User = &response.UserInfo{
					ID:         user.ID,
					Nickname:   user.Nickname,
					RealName:   user.RealName,
					Avatar:     user.Avatar,
					Department: user.Department,
				}
			}
			resp.Ranking = append(resp.Ranking, rankItem)
		}
	}

	return resp, nil
}

// StopGame 停止游戏
func (s *GameService) StopGame(roundId uint) (*response.DrawResultResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	if safeInt(round.Status) != RoundStatusPlaying {
		return nil, errors.New("游戏未在进行中")
	}

	return s.settleGame(roundId)
}

// CancelGame 取消游戏
func (s *GameService) CancelGame(roundId uint) error {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return errors.New("场次不存在")
	}

	status := safeInt(round.Status)
	if status == RoundStatusFinished {
		return errors.New("已结束的场次不能取消")
	}

	// 重置数据库状态
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":     RoundStatusPending,
		"start_time": nil,
	})

	// 清除Redis数据
	common.ClearCurrentRound(round.ActivityId)
	common.ClearRoundScores(roundId)

	// 清除数据库成绩
	global.GVA_DB.Where("round_id = ?", roundId).Delete(&annual.AnnualShakeScore{})

	return nil
}

// GetGameStatus 获取游戏状态
func (s *GameService) GetGameStatus(roundId uint) (*response.GameStatusResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	result := &response.GameStatusResp{
		RoundId:   roundId,
		Status:    safeInt(round.Status),
		Duration:  round.Duration,
		StartTime: round.StartTime,
	}

	// 从Redis获取参与人数
	playerCount, err := common.GetPlayerCount(roundId)
	if err != nil {
		// 降级查数据库
		global.GVA_DB.Model(&annual.AnnualShakeScore{}).Where("round_id = ?", roundId).Count(&playerCount)
	}
	result.PlayerCount = int(playerCount)

	// 计算剩余时间
	if safeInt(round.Status) == RoundStatusPlaying && round.StartTime != nil {
		elapsed := int(time.Since(*round.StartTime).Seconds())
		remaining := round.Duration - elapsed
		if remaining < 0 {
			remaining = 0
		}
		result.Remaining = remaining
	}

	return result, nil
}

// GetRanking 获取排行榜
func (s *GameService) GetRanking(roundId uint, limit int) (*response.RankingListResp, error) {
	if limit <= 0 {
		limit = 20
	}

	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	// 从Redis获取排行榜
	ranking, err := common.GetRanking(roundId, int64(limit))
	if err != nil {
		return nil, err
	}

	// 获取参与人数
	playerCount, _ := common.GetPlayerCount(roundId)

	list := make([]response.RankingItem, 0, len(ranking))
	for _, r := range ranking {
		item := response.RankingItem{
			Rank:     r.Rank,
			UserId:   r.UserId,
			Score:    int(r.Score),
			IsWinner: r.Rank <= round.WinnerCount,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, r.UserId).Error; err == nil {
			item.User = &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		list = append(list, item)
	}

	return &response.RankingListResp{
		List:        list,
		PlayerCount: int(playerCount),
	}, nil
}

// GetWinners 获取中奖名单
func (s *GameService) GetWinners(roundId uint) (*response.WinnerListResp, error) {
	var winners []annual.AnnualWinner
	global.GVA_DB.Where("round_id = ?", roundId).
		Order("created_at ASC").
		Find(&winners)

	list := make([]response.WinnerItem, 0, len(winners))
	for _, w := range winners {
		item := response.WinnerItem{
			ID:        w.ID,
			UserId:    w.UserId,
			WinType:   safeInt(w.WinType),
			CreatedAt: w.CreatedAt,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, w.UserId).Error; err == nil {
			item.User = &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		// 获取奖品信息
		var prize annual.AnnualPrize
		if err := global.GVA_DB.First(&prize, w.PrizeId).Error; err == nil {
			item.Prize = &response.PrizeItem{
				ID:    prize.ID,
				Name:  prize.Name,
				Image: prize.Image,
				Level: safeInt(prize.Level),
			}
		}

		// 获取分数和排名
		var score annual.AnnualShakeScore
		if err := global.GVA_DB.Where("round_id = ? AND user_id = ?", roundId, w.UserId).First(&score).Error; err == nil {
			item.Score = score.Score
			item.Rank = score.Rank
		}

		list = append(list, item)
	}

	return &response.WinnerListResp{List: list}, nil
}

// settleGame 结算游戏
func (s *GameService) settleGame(roundId uint) (*response.DrawResultResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	// 从Redis获取排行榜
	ranking, err := common.GetRanking(roundId, 0) // 获取全部
	if err != nil {
		return nil, err
	}

	// 获取奖品信息
	var prize annual.AnnualPrize
	global.GVA_DB.First(&prize, round.PrizeId)

	prizeInfo := &response.PrizeItem{
		ID:    prize.ID,
		Name:  prize.Name,
		Image: prize.Image,
		Level: safeInt(prize.Level),
	}

	winners := make([]response.WinnerItem, 0)
	winType := 1 // 摇一摇

	for _, r := range ranking {
		var user annual.AnnualUser
		global.GVA_DB.First(&user, r.UserId)

		isWinner := r.Rank <= round.WinnerCount

		// 保存成绩到数据库
		score := annual.AnnualShakeScore{
			ActivityId: round.ActivityId,
			RoundId:    roundId,
			UserId:     r.UserId,
			Score:      int(r.Score),
			Rank:       r.Rank,
			IsWinner:   boolToIntPtr(isWinner),
		}
		global.GVA_DB.Create(&score)

		// 中奖处理
		if isWinner {
			winner := annual.AnnualWinner{
				ActivityId: round.ActivityId,
				UserId:     r.UserId,
				PrizeId:    round.PrizeId,
				RoundId:    roundId,
				WinType:    &winType,
			}
			global.GVA_DB.Create(&winner)

			userInfo := &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}

			winners = append(winners, response.WinnerItem{
				ID:        winner.ID,
				Rank:      r.Rank,
				UserId:    r.UserId,
				Score:     int(r.Score),
				WinType:   winType,
				CreatedAt: winner.CreatedAt,
				User:      userInfo,
				Prize:     prizeInfo,
			})
		}
	}

	// 更新场次状态
	now := time.Now()
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":   RoundStatusFinished,
		"end_time": now,
	})

	// 更新奖品剩余数量
	if len(winners) > 0 {
		global.GVA_DB.Model(&prize).Update("remain_count", prize.RemainCount-len(winners))
	}

	// 清除Redis
	common.ClearCurrentRound(round.ActivityId)
	common.SetRoundStatus(roundId, RoundStatusFinished)

	return &response.DrawResultResp{
		Winners: winners,
		Prize:   prizeInfo,
	}, nil
}

// boolToIntPtr bool转int指针
func boolToIntPtr(b bool) *int {
	val := 0
	if b {
		val = 1
	}
	return &val
}
