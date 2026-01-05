package common

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/ws"
	"time"
)

// GetCheckInStatsWithList 获取签到统计和最新列表
func GetCheckInStatsWithList(activityId uint, limit int) *response.CheckInStatsResp {
	// 按状态统计
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

	stats := &response.CheckInStatsResp{
		List: []response.CheckInItemResp{},
	}
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

	// 获取最新签到列表
	if limit <= 0 {
		limit = 20
	}

	var checkIns []annual.AnnualCheckIn
	global.GVA_DB.Where("activity_id = ?", activityId).
		Preload("User").
		Order("check_in_time DESC").
		Limit(limit).
		Find(&checkIns)

	for _, c := range checkIns {
		stats.List = append(stats.List, response.CheckInItemResp{
			ID:          c.ID,
			RealName:    c.RealName,
			Department:  c.Department,
			CheckInTime: c.CheckInTime.Format(time.DateTime),
			Avatar:      c.User.Avatar,
			Nickname:    c.User.Nickname,
		})
	}

	return stats
}

// BroadcastCheckInStats 广播签到统计到主持人端
func BroadcastCheckInStats(activityId uint) {
	stats := GetCheckInStatsWithList(activityId, 20)
	ws.GetEventTrigger().TriggerCheckInStats(activityId, stats)
}
