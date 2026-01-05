package common

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/ws"
)

// CheckInStatsResult 签到统计结果
type CheckInStatsResult struct {
	Total    int
	Pending  int
	Approved int
	Rejected int
}

// GetCheckInStatsCount 获取签到统计数量
func GetCheckInStatsCount(activityId uint) *CheckInStatsResult {
	type StatusCount struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}

	var results []StatusCount
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Select("status, COUNT(*) as count").
		Where("activity_id = ?", activityId).
		Group("status").
		Scan(&results)

	stats := &CheckInStatsResult{}
	for _, r := range results {
		stats.Total += int(r.Count)
		switch r.Status {
		case 0:
			stats.Pending = int(r.Count)
		case 1:
			stats.Approved = int(r.Count)
		case 2:
			stats.Rejected = int(r.Count)
		}
	}

	return stats
}

// BroadcastCheckInStats 广播签到统计到主持人端
func BroadcastCheckInStats(activityId uint) {
	stats := GetCheckInStatsCount(activityId)
	ws.GetEventTrigger().TriggerCheckInStats(
		activityId,
		stats.Total,
		stats.Pending,
		stats.Approved,
		stats.Rejected,
	)
}
