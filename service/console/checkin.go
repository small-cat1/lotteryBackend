package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/service/common"
)

type CheckInService struct{}

// GetCheckInStats 获取签到统计（包含开关状态和最新列表）
func (s *CheckInService) GetCheckInStats(activityId uint, limit int) (*response.CheckInStatsResp, error) {
	// 获取签到开关状态
	isOpen, _ := common.GetCheckInSwitch(activityId)

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
		IsOpen: isOpen,
		List:   []response.CheckInItemResp{},
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
		limit = 10
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
			CheckInTime: c.CheckInTime.Format("15:04:05"),
			Avatar:      c.User.Avatar,
			Nickname:    c.User.Nickname,
		})
	}

	return stats, nil
}

// OpenCheckIn 开启签到
func (s *CheckInService) OpenCheckIn(activityId uint) error {
	// 检查活动是否存在（从缓存或数据库）
	activityService := ActivityService{}
	_, err := activityService.GetActivityDetail(activityId)
	if err != nil {
		return err
	}

	// 设置签到开关为开启
	return common.SetCheckInSwitch(activityId, true)
}

// CloseCheckIn 关闭签到
func (s *CheckInService) CloseCheckIn(activityId uint) error {
	return common.SetCheckInSwitch(activityId, false)
}

// SyncCheckInUsersToRedis 同步签到用户到Redis（启动时或数据恢复用）
func (s *CheckInService) SyncCheckInUsersToRedis(activityId uint) error {
	var userIds []uint
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ?", activityId).
		Pluck("user_id", &userIds)

	for _, userId := range userIds {
		common.AddCheckInUser(activityId, userId)
	}

	return nil
}

// IsUserCheckedIn 检查用户是否已签到
func (s *CheckInService) IsUserCheckedIn(activityId uint, userId uint) (bool, error) {
	// 先查Redis
	checked, err := common.IsUserCheckIn(activityId, userId)
	if err == nil {
		return checked, nil
	}

	// Redis出错，降级查数据库
	var count int64
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ? AND user_id = ?", activityId, userId).
		Count(&count)
	return count > 0, nil
}

// DoCheckIn 执行签到（供H5端调用）
func (s *CheckInService) DoCheckIn(activityId uint, userId uint, ip string) error {
	// 检查签到是否开启
	isOpen, _ := common.GetCheckInSwitch(activityId)
	if !isOpen {
		return errors.New("签到未开启")
	}

	// 检查是否已签到
	checked, _ := s.IsUserCheckedIn(activityId, userId)
	if checked {
		return errors.New("您已签到")
	}

	// 添加到Redis
	added, err := common.AddCheckInUser(activityId, userId)
	if err != nil {
		return errors.New("签到失败")
	}
	if !added {
		return errors.New("您已签到")
	}

	// 写入数据库
	checkIn := annual.AnnualCheckIn{
		ActivityId: activityId,
		UserId:     userId,
		Ip:         ip,
	}
	if err := global.GVA_DB.Create(&checkIn).Error; err != nil {
		return errors.New("签到失败")
	}

	return nil
}
