package annual

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"time"
)

type AnnualActivityService struct{}

// CreateActivity 创建活动
func (s *AnnualActivityService) CreateActivity(activity annual.AnnualActivity) (err error) {
	return global.GVA_DB.Create(&activity).Error
}

// GetActivityList 获取活动列表
func (s *AnnualActivityService) GetActivityList(info annualReq.ActivitySearch) (list []annual.AnnualActivity, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualActivity{})

	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id DESC").Find(&list).Error
	return
}

// GetActivityById 获取活动详情
func (s *AnnualActivityService) GetActivityById(id string) (activity annual.AnnualActivity, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&activity).Error
	return
}

// UpdateActivity 更新活动
func (s *AnnualActivityService) UpdateActivity(activity annual.AnnualActivity) (err error) {
	return global.GVA_DB.Model(&annual.AnnualActivity{}).Where("id = ?", activity.ID).Updates(&activity).Error
}

// UpdateActivityStatus 更新活动状态
func (s *AnnualActivityService) UpdateActivityStatus(id uint, status int) (err error) {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == 1 {
		updates["start_time"] = time.Now()
		updates["end_time"] = time.Now().Add(2 * time.Hour)
	}
	if status == 2 {
		updates["end_time"] = time.Now()
	}
	return global.GVA_DB.Model(&annual.AnnualActivity{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteActivity 删除活动（安全检查）
func (s *AnnualActivityService) DeleteActivity(id uint) (err error) {
	// 检查关联数据
	if err = s.checkActivityRelations(id); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&annual.AnnualActivity{}, id).Error
}

// DeleteActivityByIds 批量删除活动（安全检查）
func (s *AnnualActivityService) DeleteActivityByIds(ids []uint) (err error) {
	for _, id := range ids {
		if err = s.checkActivityRelations(id); err != nil {
			return err
		}
	}
	return global.GVA_DB.Delete(&annual.AnnualActivity{}, ids).Error
}

// checkActivityRelations 检查活动关联数据
func (s *AnnualActivityService) checkActivityRelations(activityId uint) error {
	var count int64

	// 检查签到记录
	if err := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该活动下存在签到记录，无法删除")
	}

	// 检查弹幕记录
	if err := global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("activity_id = ?", activityId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该活动下存在弹幕记录，无法删除")
	}

	// 检查奖品
	if err := global.GVA_DB.Model(&annual.AnnualPrize{}).Where("activity_id = ?", activityId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该活动下存在奖品，无法删除")
	}

	// 检查摇一摇场次
	if err := global.GVA_DB.Model(&annual.AnnualShakeRound{}).Where("activity_id = ?", activityId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该活动下存在摇一摇场次，无法删除")
	}

	// 检查中奖记录
	if err := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("activity_id = ?", activityId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该活动下存在中奖记录，无法删除")
	}

	return nil
}
