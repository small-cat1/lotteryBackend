package console

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/pkg/cache"
	"lotteryBackend/ws"
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
// StartGame 开始游戏（验证密码）
// 修改：返回 endTime
func (s *GameService) StartGame(roundId uint, password string) (int64, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.Preload("Prize").First(&round, roundId).Error; err != nil {
		return 0, errors.New("场次不存在")
	}

	// 验证密码
	if round.Password != password {
		return 0, errors.New("密码错误")
	}

	// 检查是否有进行中的游戏
	currentRoundId, err := cache.GetCurrentRound(round.ActivityId)
	if err == nil && currentRoundId > 0 {
		return 0, errors.New("当前有进行中的游戏，请先结束后再开始")
	}

	// 计算时间
	now := time.Now()
	duration := round.Duration
	if duration <= 0 {
		duration = 60
	}
	endTime := now.Add(time.Duration(duration) * time.Second)
	endTimeMs := endTime.UnixMilli()

	// 更新数据库状态
	global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":     RoundStatusPlaying,
		"start_time": now,
	})

	// 设置Redis
	cache.SetCurrentRound(round.ActivityId, roundId, duration)
	cache.SetRoundStatus(roundId, RoundStatusPlaying)

	// ⭐ 广播游戏开始
	gameStartPayload := ws.GameStartPayload{
		RoundId:  roundId,
		Duration: duration,
		EndTime:  endTimeMs,
		Round: ws.RoundInfo{
			ID:          round.ID,
			RoundName:   round.RoundName,
			Duration:    round.Duration,
			WinnerCount: round.WinnerCount,
			Status:      RoundStatusPlaying,
		},
	}

	if round.Prize.ID > 0 {
		gameStartPayload.Round.Prize = ws.PrizeBrief{
			ID:    round.Prize.ID,
			Name:  round.Prize.Name,
			Image: round.Prize.Image,
			Level: safeInt(round.Prize.Level),
		}
	}

	ws.GetEventTrigger().TriggerGameStart(round.ActivityId, gameStartPayload)

	return endTimeMs, nil
}

// GetCurrent 获取当前进行中的游戏
// GetCurrent 获取当前进行中的游戏
// 修改：返回 endTime，移除 remaining
func (s *GameService) GetCurrent(activityId uint) (*response.CurrentGameResp, error) {
	resp := &response.CurrentGameResp{
		Round:   nil,
		Status:  0,
		EndTime: 0,
		Ranking: []response.RankingItem{},
	}

	// 获取当前场次ID
	currentRoundId, err := cache.GetCurrentRound(activityId)
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

	resp.Status = RoundStatusPlaying

	// ⭐ 计算 endTime（从数据库的 start_time + duration）
	if round.StartTime != nil {
		endTime := round.StartTime.Add(time.Duration(round.Duration) * time.Second)
		resp.EndTime = endTime.UnixMilli()
	}

	// 获取参与人数
	playerCount, _ := cache.GetPlayerCount(currentRoundId)
	resp.PlayerCount = int(playerCount)

	// 获取排行榜
	ranking, _ := cache.GetRanking(currentRoundId, 20)
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
					ID:       user.ID,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
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

	// 第一步：先关闭 Redis 游戏状态
	cache.ClearCurrentRound(round.ActivityId)
	cache.SetRoundStatus(roundId, RoundStatusFinished)
	// ⭐ 广播游戏结束
	ws.GetEventTrigger().TriggerGameStop(round.ActivityId, roundId)

	result, err := s.settleGame(roundId)
	if err != nil {
		global.GVA_LOG.Error("StopGame 结算失败", zap.Uint("roundId", roundId), zap.Error(err))
		return nil, err
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
	ranking, err := cache.GetRanking(roundId, int64(limit))
	if err != nil {
		return nil, err
	}

	// 获取参与人数
	playerCount, _ := cache.GetPlayerCount(roundId)

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
				ID:       user.ID,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
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
				ID:       user.ID,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
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
// settleGame 结算游戏
func (s *GameService) settleGame(roundId uint) (*response.DrawResultResp, error) {

	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		global.GVA_LOG.Error("settleGame 场次不存在", zap.Uint("roundId", roundId), zap.Error(err))
		return nil, errors.New("场次不存在")
	}

	// 从Redis获取排行榜
	ranking, err := cache.GetRanking(roundId, 0)
	if err != nil {
		global.GVA_LOG.Error("settleGame 获取排行榜失败", zap.Uint("roundId", roundId), zap.Error(err))
		return nil, err
	}

	if len(ranking) == 0 {
		// ✅ 仍然需要更新数据库状态为已结束
		if err := global.GVA_DB.Model(&round).Updates(map[string]interface{}{
			"status":   RoundStatusFinished,
			"end_time": time.Now(),
		}).Error; err != nil {
			global.GVA_LOG.Error("settleGame 更新场次状态失败", zap.Error(err))
			return nil, fmt.Errorf("更新场次状态失败: %w", err)
		}
		return &response.DrawResultResp{
			Winners: []response.WinnerItem{},
			Prize:   nil,
		}, nil
	}

	// 获取奖品信息
	var prize annual.AnnualPrize
	if err := global.GVA_DB.First(&prize, round.PrizeId).Error; err != nil {
		global.GVA_LOG.Error("settleGame 奖品不存在", zap.Uint("prizeId", round.PrizeId), zap.Error(err))
		return nil, errors.New("奖品不存在")
	}

	prizeInfo := &response.PrizeItem{
		ID:    prize.ID,
		Name:  prize.Name,
		Image: prize.Image,
		Level: safeInt(prize.Level),
	}

	// 提前批量获取所有用户信息
	userIds := make([]uint, 0, len(ranking))
	for _, r := range ranking {
		userIds = append(userIds, r.UserId)
	}

	var users []annual.AnnualUser
	if err := global.GVA_DB.Where("id IN ?", userIds).Find(&users).Error; err != nil {
		global.GVA_LOG.Error("settleGame 获取用户信息失败", zap.Error(err))
		return nil, errors.New("获取用户信息失败")
	}

	// 构建用户Map
	userMap := make(map[uint]annual.AnnualUser, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 准备批量插入的数据
	scores := make([]annual.AnnualShakeScore, 0, len(ranking))
	winnersData := make([]annual.AnnualWinner, 0, round.WinnerCount)
	winners := make([]response.WinnerItem, 0, round.WinnerCount)
	winType := 1

	for _, r := range ranking {
		isWinner := r.Rank <= round.WinnerCount

		scores = append(scores, annual.AnnualShakeScore{
			ActivityId: round.ActivityId,
			RoundId:    roundId,
			UserId:     r.UserId,
			Score:      int(r.Score),
			Rank:       r.Rank,
			IsWinner:   boolToIntPtr(isWinner),
		})

		if isWinner {
			winnersData = append(winnersData, annual.AnnualWinner{
				ActivityId: round.ActivityId,
				UserId:     r.UserId,
				PrizeId:    round.PrizeId,
				RoundId:    roundId,
				WinType:    &winType,
			})

			user := userMap[r.UserId]
			userInfo := &response.UserInfo{
				ID:       user.ID,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
			}

			winners = append(winners, response.WinnerItem{
				Rank:    r.Rank,
				UserId:  r.UserId,
				Score:   int(r.Score),
				WinType: winType,
				User:    userInfo,
				Prize:   prizeInfo,
			})
		}
	}

	// 使用事务批量写入
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if len(scores) > 0 {
			if err := tx.CreateInBatches(&scores, 100).Error; err != nil {
				global.GVA_LOG.Error("settleGame 保存成绩失败", zap.Error(err))
				return fmt.Errorf("保存成绩失败: %w", err)
			}
		}

		if len(winnersData) > 0 {
			if err := tx.CreateInBatches(&winnersData, 100).Error; err != nil {
				global.GVA_LOG.Error("settleGame 保存中奖记录失败", zap.Error(err))
				return fmt.Errorf("保存中奖记录失败: %w", err)
			}

			result := tx.Model(&prize).
				Where("remain_count >= ?", len(winnersData)).
				Update("remain_count", gorm.Expr("remain_count - ?", len(winnersData)))

			if result.Error != nil {
				global.GVA_LOG.Error("settleGame 更新奖品数量失败", zap.Error(result.Error))
				return fmt.Errorf("更新奖品数量失败: %w", result.Error)
			}

			// 如果没有更新到行，说明库存不足，但不影响结算
			if result.RowsAffected == 0 {
				global.GVA_LOG.Warn("settleGame 奖品库存不足，跳过扣减", zap.Uint("prizeId", prize.ID))
			}
		}

		if err := tx.Model(&round).Updates(map[string]interface{}{
			"status":   RoundStatusFinished,
			"end_time": time.Now(),
		}).Error; err != nil {
			global.GVA_LOG.Error("settleGame 更新场次状态失败", zap.Error(err))
			return fmt.Errorf("更新场次状态失败: %w", err)
		}
		return nil
	})

	if err != nil {
		global.GVA_LOG.Error("settleGame 事务失败", zap.Uint("roundId", roundId), zap.Error(err))
		return nil, err
	}

	// 回填中奖记录ID
	for i, w := range winnersData {
		winners[i].ID = w.ID
		winners[i].CreatedAt = w.CreatedAt
	}

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
