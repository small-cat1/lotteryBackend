package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
	"time"
)

type H5CheckInService struct{}

// CheckIn 用户签到
func (s *H5CheckInService) CheckIn(userId, activityId uint, ip string) (*response.CheckInRecordResp, error) {
	// 检查活动是否存在且进行中
	var activity annual.AnnualActivity
	if err := global.GVA_DB.First(&activity, activityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	if *activity.Status != 1 {
		return nil, errors.New("活动未开始或已结束")
	}

	if *activity.CheckInEnabled != 1 {
		return nil, errors.New("签到功能未开启")
	}

	// 检查是否已签到
	var existCheckIn annual.AnnualCheckIn
	result := global.GVA_DB.Where("activity_id = ? AND user_id = ?", activityId, userId).First(&existCheckIn)
	if result.RowsAffected > 0 {
		return nil, errors.New("您已签到，请勿重复签到")
	}

	// 创建签到记录
	checkIn := annual.AnnualCheckIn{
		ActivityId:  activityId,
		UserId:      userId,
		CheckInTime: time.Now(),
		Ip:          ip,
	}

	if err := global.GVA_DB.Create(&checkIn).Error; err != nil {
		return nil, errors.New("签到失败，请稍后重试")
	}

	// 获取用户信息
	userService := H5UserService{}
	userBrief := userService.GetUserBrief(userId)

	return &response.CheckInRecordResp{
		ID:          checkIn.ID,
		User:        *userBrief,
		CheckInTime: checkIn.CheckInTime,
	}, nil
}

// GetCheckInStatus 获取签到状态
func (s *H5CheckInService) GetCheckInStatus(userId, activityId uint) (*response.CheckInStatusResp, error) {
	var checkIn annual.AnnualCheckIn
	result := global.GVA_DB.Where("activity_id = ? AND user_id = ?", activityId, userId).First(&checkIn)

	if result.RowsAffected == 0 {
		return &response.CheckInStatusResp{
			IsCheckedIn: false,
			CheckInTime: nil,
		}, nil
	}

	return &response.CheckInStatusResp{
		IsCheckedIn: true,
		CheckInTime: &checkIn.CheckInTime,
	}, nil
}

// GetCheckInList 获取签到列表
func (s *H5CheckInService) GetCheckInList(activityId uint, page, pageSize int) (*response.H5PageResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var checkIns []annual.AnnualCheckIn
	var total int64

	db := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId)
	db.Count(&total)

	err := db.Order("check_in_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&checkIns).Error

	if err != nil {
		return nil, err
	}

	// 转换响应
	userService := H5UserService{}
	list := make([]response.CheckInRecordResp, len(checkIns))
	for i, c := range checkIns {
		userBrief := userService.GetUserBrief(c.UserId)
		list[i] = response.CheckInRecordResp{
			ID:          c.ID,
			CheckInTime: c.CheckInTime,
		}
		if userBrief != nil {
			list[i].User = *userBrief
		}
	}

	return &response.H5PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetCheckInStats 获取签到统计
func (s *H5CheckInService) GetCheckInStats(activityId uint) (*response.CheckInStatsResp, error) {
	var checkedCount int64
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId).Count(&checkedCount)

	// 获取已报名通过的总人数
	var totalCount int64
	global.GVA_DB.Model(&annual.AnnualUser{}).Where("status = ?", 1).Count(&totalCount)

	var checkRate float64 = 0
	if totalCount > 0 {
		checkRate = float64(checkedCount) / float64(totalCount) * 100
	}

	return &response.CheckInStatsResp{
		CheckedCount: int(checkedCount),
		TotalCount:   int(totalCount),
		CheckRate:    checkRate,
	}, nil
}

// GetRecentCheckIns 获取最新签到
func (s *H5CheckInService) GetRecentCheckIns(activityId uint, limit int) ([]response.CheckInRecordResp, error) {
	if limit <= 0 {
		limit = 20
	}

	var checkIns []annual.AnnualCheckIn
	err := global.GVA_DB.Where("activity_id = ?", activityId).
		Order("check_in_time DESC").
		Limit(limit).
		Find(&checkIns).Error

	if err != nil {
		return nil, err
	}

	userService := H5UserService{}
	result := make([]response.CheckInRecordResp, len(checkIns))
	for i, c := range checkIns {
		userBrief := userService.GetUserBrief(c.UserId)
		result[i] = response.CheckInRecordResp{
			ID:          c.ID,
			CheckInTime: c.CheckInTime,
		}
		if userBrief != nil {
			result[i].User = *userBrief
		}
	}

	return result, nil
}
