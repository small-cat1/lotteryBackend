package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
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
	return global.GVA_DB.Model(&annual.AnnualActivity{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteActivity 删除活动
func (s *AnnualActivityService) DeleteActivity(id uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualActivity{}, id).Error
}

// DeleteActivityByIds 批量删除活动
func (s *AnnualActivityService) DeleteActivityByIds(ids []uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualActivity{}, ids).Error
}
