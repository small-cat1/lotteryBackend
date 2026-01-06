package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
)

type H5ActivityService struct{}

// GetActivityDetail 获取活动详情
func (s *H5ActivityService) GetActivityDetail(activityId uint) (*response.H5ActivityResp, error) {
	var activity annual.AnnualActivity

	if err := global.GVA_DB.First(&activity, activityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	return s.toActivityResp(&activity), nil
}

// CheckActivityOngoing 检查活动是否进行中
func (s *H5ActivityService) CheckActivityOngoing(activityId uint) error {
	var activity annual.AnnualActivity

	if err := global.GVA_DB.First(&activity, activityId).Error; err != nil {
		return errors.New("活动不存在")
	}

	if *activity.Status != 1 {
		return errors.New("活动未开始或已结束")
	}

	return nil
}

// toActivityResp 转换为响应结构
func (s *H5ActivityService) toActivityResp(activity *annual.AnnualActivity) *response.H5ActivityResp {
	return &response.H5ActivityResp{
		ID:             activity.ID,
		Title:          activity.Title,
		Cover:          activity.Cover,
		Description:    activity.Description,
		StartTime:      *activity.StartTime,
		EndTime:        *activity.EndTime,
		Status:         *activity.Status,
		CheckInEnabled: *activity.CheckInEnabled,
		DanmakuEnabled: *activity.DanmakuEnabled,
		DanmakuAudit:   *activity.DanmakuAudit,
	}
}
