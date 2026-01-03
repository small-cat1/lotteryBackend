package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualCheckInService struct{}

// GetCheckInList 获取签到列表
func (s *AnnualCheckInService) GetCheckInList(info annualReq.CheckInSearch) (list []annual.AnnualCheckIn, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Preload("User")

	if info.ActivityId != 0 {
		db = db.Where("activity_id = ?", info.ActivityId)
	}
	if info.UserId != 0 {
		db = db.Where("user_id = ?", info.UserId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id DESC").Find(&list).Error
	return
}

// GetCheckInStats 签到统计
func (s *AnnualCheckInService) GetCheckInStats(activityId string) (stats map[string]interface{}, err error) {
	var total int64
	var registered int64

	// 总签到人数
	err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId).Count(&total).Error
	if err != nil {
		return
	}

	// 已报名签到人数
	err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Joins("LEFT JOIN annual_users ON annual_check_ins.user_id = annual_users.id").
		Where("annual_check_ins.activity_id = ? AND annual_users.is_registered = 1", activityId).
		Count(&registered).Error
	if err != nil {
		return
	}

	stats = map[string]interface{}{
		"total":      total,
		"registered": registered,
		"guest":      total - registered,
	}
	return
}

// ExportCheckIn 导出签到
func (s *AnnualCheckInService) ExportCheckIn(info annualReq.CheckInSearch) (filePath string, err error) {
	// TODO: 实现导出逻辑
	return "", nil
}
