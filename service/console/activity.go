package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/pkg/cache"
)

type ActivityService struct{}

// GetActivityDetail 获取活动详情（先查缓存，没有再查库）
func (s *ActivityService) GetActivityDetail(activityId uint) (*response.ActivityDetailResp, error) {
	// 先查缓存
	activityCache, err := cache.GetActivityCache(activityId)
	if err == nil && activityCache != nil {
		return &response.ActivityDetailResp{
			ID:             activityCache.ID,
			Title:          activityCache.Title,
			Logo:           activityCache.Logo,
			Cover:          activityCache.Cover,
			Description:    activityCache.Description,
			CheckInEnabled: activityCache.CheckInEnabled,
			DanmakuEnabled: activityCache.DanmakuEnabled,
			DanmakuAudit:   activityCache.DanmakuAudit,
			WinnerExclude:  activityCache.WinnerExclude,
			Status:         activityCache.Status,
		}, nil
	}

	// 缓存没有，查数据库
	var activity annual.AnnualActivity
	if err := global.GVA_DB.First(&activity, activityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	// 构建响应
	resp := &response.ActivityDetailResp{
		ID:             activity.ID,
		Title:          activity.Title,
		Logo:           "",
		Cover:          activity.Cover,
		Description:    activity.Description,
		StartTime:      activity.StartTime,
		EndTime:        activity.EndTime,
		CheckInEnabled: safeInt(activity.CheckInEnabled),
		DanmakuEnabled: safeInt(activity.DanmakuEnabled),
		DanmakuAudit:   safeInt(activity.DanmakuAudit),
		WinnerExclude:  safeInt(activity.WinnerExclude),
		Status:         safeInt(activity.Status),
	}

	// 存入缓存
	cache.SetActivityCache(&cache.ActivityCache{
		ID:             activity.ID,
		Title:          activity.Title,
		Logo:           "",
		Cover:          activity.Cover,
		Description:    activity.Description,
		CheckInEnabled: safeInt(activity.CheckInEnabled),
		DanmakuEnabled: safeInt(activity.DanmakuEnabled),
		DanmakuAudit:   safeInt(activity.DanmakuAudit),
		WinnerExclude:  safeInt(activity.WinnerExclude),
		Status:         safeInt(activity.Status),
	})

	return resp, nil
}

// GetPrizeList 获取奖品列表
func (s *ActivityService) GetPrizeList(activityId uint) (*response.PrizeListResp, error) {
	var prizes []annual.AnnualPrize
	global.GVA_DB.Where("activity_id = ?", activityId).
		Order("level ASC, id ASC").
		Find(&prizes)

	list := make([]response.PrizeItem, 0, len(prizes))
	for _, p := range prizes {
		list = append(list, response.PrizeItem{
			ID:          p.ID,
			Name:        p.Name,
			Image:       p.Image,
			Level:       safeInt(p.Level),
			TotalCount:  p.TotalCount,
			RemainCount: p.RemainCount,
		})
	}

	return &response.PrizeListResp{List: list}, nil
}

// GetRoundList 获取场次列表
func (s *ActivityService) GetRoundList(activityId uint) (*response.RoundListResp, error) {
	var rounds []annual.AnnualShakeRound
	global.GVA_DB.Where("activity_id = ?", activityId).
		Order("sort ASC, id ASC").
		Find(&rounds)

	list := make([]response.RoundItem, 0, len(rounds))
	for _, r := range rounds {
		item := response.RoundItem{
			ID:          r.ID,
			RoundName:   r.RoundName,
			Duration:    r.Duration,
			WinnerCount: r.WinnerCount,
			Status:      safeInt(r.Status),
		}

		// 获取奖品信息
		if r.PrizeId > 0 {
			var prize annual.AnnualPrize
			if err := global.GVA_DB.First(&prize, r.PrizeId).Error; err == nil {
				item.Prize = &response.PrizeItem{
					ID:    prize.ID,
					Name:  prize.Name,
					Image: prize.Image,
					Level: safeInt(prize.Level),
				}
			}
		}

		// 统计实际中奖人数
		var winnerCount int64
		global.GVA_DB.Model(&annual.AnnualWinner{}).
			Where("round_id = ?", r.ID).
			Count(&winnerCount)
		item.ActualWinners = int(winnerCount)

		list = append(list, item)
	}

	return &response.RoundListResp{List: list}, nil
}

// GetRoundDetail 获取场次详情
func (s *ActivityService) GetRoundDetail(roundId uint) (*response.RoundDetailResp, error) {
	var round annual.AnnualShakeRound
	if err := global.GVA_DB.First(&round, roundId).Error; err != nil {
		return nil, errors.New("场次不存在")
	}

	result := &response.RoundDetailResp{
		ID:          round.ID,
		ActivityId:  round.ActivityId,
		RoundName:   round.RoundName,
		Duration:    round.Duration,
		WinnerCount: round.WinnerCount,
		Status:      safeInt(round.Status),
		StartTime:   round.StartTime,
		EndTime:     round.EndTime,
	}

	// 获取奖品信息
	if round.PrizeId > 0 {
		var prize annual.AnnualPrize
		if err := global.GVA_DB.First(&prize, round.PrizeId).Error; err == nil {
			result.Prize = &response.PrizeItem{
				ID:    prize.ID,
				Name:  prize.Name,
				Image: prize.Image,
				Level: safeInt(prize.Level),
			}
		}
	}

	// 统计
	var winnerCount int64
	global.GVA_DB.Model(&annual.AnnualWinner{}).Where("round_id = ?", roundId).Count(&winnerCount)
	result.ActualWinners = int(winnerCount)

	var playerCount int64
	global.GVA_DB.Model(&annual.AnnualShakeScore{}).Where("round_id = ?", roundId).Count(&playerCount)
	result.PlayerCount = int(playerCount)

	return result, nil
}
