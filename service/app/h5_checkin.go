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
