package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
	"time"
)

type H5ActivityService struct{}

// GetCurrentActivity 获取当前进行中的活动
func (s *H5ActivityService) GetCurrentActivity() (*response.H5ActivityResp, error) {
	var activity annual.AnnualActivity
	now := time.Now()

	// 优先获取进行中的活动
	result := global.GVA_DB.Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).
		Order("created_at DESC").First(&activity)

	if result.RowsAffected == 0 {
		// 获取最近一个活动
		result = global.GVA_DB.Order("created_at DESC").First(&activity)
		if result.RowsAffected == 0 {
			return nil, errors.New("暂无活动")
		}
	}

	return s.toActivityResp(&activity), nil
}

// GetActivityList 获取活动列表
func (s *H5ActivityService) GetActivityList() ([]response.H5ActivityResp, error) {
	var activities []annual.AnnualActivity

	err := global.GVA_DB.Order("created_at DESC").Find(&activities).Error
	if err != nil {
		return nil, err
	}

	result := make([]response.H5ActivityResp, len(activities))
	for i, activity := range activities {
		result[i] = *s.toActivityResp(&activity)
	}

	return result, nil
}

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
