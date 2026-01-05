package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
)

type AnnualDashboardService struct{}

// DashboardStats 统计概览数据
type DashboardStats struct {
	CheckInCount        int64 `json:"checkInCount"`
	RegisteredCount     int64 `json:"registeredCount"`
	DanmakuCount        int64 `json:"danmakuCount"`
	PendingDanmakuCount int64 `json:"pendingDanmakuCount"`
	ShakeCount          int64 `json:"shakeCount"`
	ShakeRoundCount     int64 `json:"shakeRoundCount"`
	WinnerCount         int64 `json:"winnerCount"`
	ReceivedCount       int64 `json:"receivedCount"`
}

// CheckInTrendItem 签到趋势项
type CheckInTrendItem struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

// PrizeStatsItem 奖品统计项
type PrizeStatsItem struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
	Used  int    `json:"used"`
}

// RecentWinnerItem 最新中奖项
type RecentWinnerItem struct {
	ID        uint   `json:"id"`
	UserName  string `json:"userName"`
	Avatar    string `json:"avatar"`
	PrizeName string `json:"prizeName"`
	Time      string `json:"time"`
}

// HotWordItem 热词项
type HotWordItem struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

// RecentDanmakuItem 最新弹幕项
type RecentDanmakuItem struct {
	ID       uint   `json:"id"`
	UserName string `json:"userName"`
	Content  string `json:"content"`
	Color    string `json:"color"`
}

// ShakeRankingItem 摇一摇排行项
type ShakeRankingItem struct {
	Rank     int    `json:"rank"`
	UserName string `json:"userName"`
	Score    int    `json:"score"`
}

// GetDashboardStats 获取统计概览数据
func (s *AnnualDashboardService) GetDashboardStats(activityId string) (stats DashboardStats, err error) {
	db := global.GVA_DB

	// 签到人数
	db.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId).Count(&stats.CheckInCount)

	// 已报名签到人数
	db.Model(&annual.AnnualCheckIn{}).
		Joins("LEFT JOIN annual_users ON annual_check_ins.user_id = annual_users.id").
		Where("annual_check_ins.activity_id = ? AND annual_users.is_registered = 1", activityId).
		Count(&stats.RegisteredCount)

	// 弹幕总数
	db.Model(&annual.AnnualDanmaku{}).Where("activity_id = ?", activityId).Count(&stats.DanmakuCount)

	// 待审核弹幕数
	db.Model(&annual.AnnualDanmaku{}).Where("activity_id = ? AND status = 0", activityId).Count(&stats.PendingDanmakuCount)

	// 摇一摇总次数
	db.Model(&annual.AnnualShakeScore{}).
		Joins("LEFT JOIN annual_shake_rounds ON annual_shake_scores.round_id = annual_shake_rounds.id").
		Where("annual_shake_rounds.activity_id = ?", activityId).
		Select("COALESCE(SUM(annual_shake_scores.score), 0)").
		Scan(&stats.ShakeCount)

	// 摇一摇场次数
	db.Model(&annual.AnnualShakeRound{}).Where("activity_id = ?", activityId).Count(&stats.ShakeRoundCount)

	// 中奖人数
	db.Model(&annual.AnnualWinner{}).Where("activity_id = ?", activityId).Count(&stats.WinnerCount)

	// 已领奖人数
	db.Model(&annual.AnnualWinner{}).Where("activity_id = ? AND status = 1", activityId).Count(&stats.ReceivedCount)

	return
}

// GetCheckInTrend 获取签到趋势
func (s *AnnualDashboardService) GetCheckInTrend(activityId string) (list []CheckInTrendItem, err error) {
	err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Select("DATE_FORMAT(check_in_time, '%H:%i') as time, COUNT(*) as count").
		Where("activity_id = ?", activityId).
		Group("DATE_FORMAT(check_in_time, '%H:%i')").
		Order("time ASC").
		Find(&list).Error
	return
}

// GetPrizeStats 获取奖品统计
func (s *AnnualDashboardService) GetPrizeStats(activityId string) (list []PrizeStatsItem, err error) {
	var prizes []annual.AnnualPrize
	err = global.GVA_DB.Where("activity_id = ?", activityId).Order("level ASC").Find(&prizes).Error
	if err != nil {
		return
	}

	for _, prize := range prizes {
		list = append(list, PrizeStatsItem{
			Name:  prize.Name,
			Total: prize.TotalCount,
			Used:  prize.TotalCount - prize.RemainCount,
		})
	}
	return
}

// GetRecentWinners 获取最新中奖记录
func (s *AnnualDashboardService) GetRecentWinners(activityId string, limit int) (list []RecentWinnerItem, err error) {
	var winners []annual.AnnualWinner
	err = global.GVA_DB.
		Preload("User").
		Preload("Prize").
		Where("activity_id = ?", activityId).
		Order("created_at DESC").
		Limit(limit).
		Find(&winners).Error
	if err != nil {
		return
	}

	for _, winner := range winners {
		userName := ""
		avatar := ""
		userName = winner.User.Nickname
		avatar = winner.User.Avatar

		prizeName := ""
		if winner.Prize.Name != "" {
			prizeName = winner.Prize.Name
		}

		list = append(list, RecentWinnerItem{
			ID:        winner.ID,
			UserName:  userName,
			Avatar:    avatar,
			PrizeName: prizeName,
			Time:      winner.CreatedAt.Format("15:04:05"),
		})
	}
	return
}

// GetHotWords 获取弹幕热词（简化版，统计关键词出现次数）
func (s *AnnualDashboardService) GetHotWords(activityId string, limit int) (list []HotWordItem, err error) {
	// 这里使用简化的实现，实际项目中可以使用分词库
	// 这里我们返回一些模拟数据，实际应用中需要做分词处理

	// 统计弹幕中常见的关键词
	var danmakus []annual.AnnualDanmaku
	err = global.GVA_DB.Where("activity_id = ? AND status = 1", activityId).
		Order("created_at DESC").
		Limit(500).
		Find(&danmakus).Error
	if err != nil {
		return
	}

	// 简单的关键词统计（实际项目中应该使用分词库）
	wordCount := make(map[string]int64)
	keywords := []string{"新年快乐", "2026", "加油", "发财", "冲冲冲", "666", "感谢", "开心", "福利", "牛"}

	for _, danmaku := range danmakus {
		for _, keyword := range keywords {
			if contains(danmaku.Content, keyword) {
				wordCount[keyword]++
			}
		}
	}

	for word, count := range wordCount {
		if count > 0 {
			list = append(list, HotWordItem{Name: word, Value: count})
		}
	}

	// 按数量排序并限制返回数量
	// 这里简化处理，实际应该排序
	if len(list) > limit {
		list = list[:limit]
	}

	return
}

// GetRecentDanmaku 获取最新弹幕
func (s *AnnualDashboardService) GetRecentDanmaku(activityId string, limit int) (list []RecentDanmakuItem, err error) {
	var danmakus []annual.AnnualDanmaku
	err = global.GVA_DB.
		Preload("User").
		Where("activity_id = ? AND status = 1", activityId).
		Order("created_at DESC").
		Limit(limit).
		Find(&danmakus).Error
	if err != nil {
		return
	}

	for _, danmaku := range danmakus {
		userName := ""
		userName = danmaku.User.Nickname
		list = append(list, RecentDanmakuItem{
			ID:       danmaku.ID,
			UserName: userName,
			Content:  danmaku.Content,
			Color:    danmaku.Color,
		})
	}
	return
}

// GetShakeRanking 获取摇一摇排行榜（当前进行中或最近一场的排行）
func (s *AnnualDashboardService) GetShakeRanking(activityId string, limit int) (list []ShakeRankingItem, err error) {
	// 找到当前进行中或最近结束的场次
	var round annual.AnnualShakeRound
	err = global.GVA_DB.Where("activity_id = ?", activityId).
		Order("status DESC, id DESC").
		First(&round).Error
	if err != nil {
		return []ShakeRankingItem{}, nil
	}

	var scores []annual.AnnualShakeScore
	err = global.GVA_DB.
		Preload("User").
		Where("round_id = ?", round.ID).
		Order("score DESC").
		Limit(limit).
		Find(&scores).Error
	if err != nil {
		return
	}

	for i, score := range scores {
		userName := ""
		userName = score.User.Nickname
		list = append(list, ShakeRankingItem{
			Rank:     i + 1,
			UserName: userName,
			Score:    score.Score,
		})
	}
	return
}

// 简单的字符串包含判断
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
