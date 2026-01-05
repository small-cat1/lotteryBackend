package annual

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
	"lotteryBackend/service/common"
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
	if info.RealName != "" {
		db = db.Where("real_name LIKE ?", "%"+info.RealName+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.Department != "" {
		db = db.Where("department LIKE ?", "%"+info.Department+"%")
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

// GetCheckInStats 签到统计
func (s *AnnualCheckInService) GetCheckInStats(activityId string) (stats map[string]interface{}, err error) {
	var total, pending, approved, rejected int64

	db := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId)

	// 总签到人数
	if err = db.Count(&total).Error; err != nil {
		return
	}

	// 待审核
	if err = db.Where("status = ?", 0).Count(&pending).Error; err != nil {
		return
	}

	// 已通过
	if err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ? AND status = ?", activityId, 1).Count(&approved).Error; err != nil {
		return
	}

	// 已拒绝
	if err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ? AND status = ?", activityId, 2).Count(&rejected).Error; err != nil {
		return
	}

	stats = map[string]interface{}{
		"total":    total,
		"pending":  pending,
		"approved": approved,
		"rejected": rejected,
	}
	return
}

// UpdateCheckIn 更新签到信息
func (s *AnnualCheckInService) UpdateCheckIn(req annualReq.CheckInUpdate) error {
	return global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("id = ?", req.Id).Updates(map[string]interface{}{
		"real_name":     req.RealName,
		"phone":         req.Phone,
		"department":    req.Department,
		"employee_no":   req.EmployeeNo,
		"status":        req.Status,
		"reject_reason": req.RejectReason,
	}).Error
}

// UpdateCheckInStatus 更新签到状态（支持批量）
func (s *AnnualCheckInService) UpdateCheckInStatus(req annualReq.CheckInStatusUpdate) error {
	updates := map[string]interface{}{
		"status": req.Status,
	}
	if req.Status == 2 {
		updates["reject_reason"] = req.RejectReason
	}

	var err error
	var activityId uint

	// 批量审核
	if len(req.Ids) > 0 {
		// 先获取 activityId
		var checkIn annual.AnnualCheckIn
		global.GVA_DB.Select("activity_id").First(&checkIn, req.Ids[0])
		activityId = checkIn.ActivityId

		err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
			Where("id IN ? AND status = 0", req.Ids).
			Updates(updates).Error
	} else if req.Id > 0 {
		// 单个审核 - 先获取 activityId
		var checkIn annual.AnnualCheckIn
		global.GVA_DB.Select("activity_id").First(&checkIn, req.Id)
		activityId = checkIn.ActivityId

		err = global.GVA_DB.Model(&annual.AnnualCheckIn{}).
			Where("id = ? AND status = 0", req.Id).
			Updates(updates).Error
	} else {
		return errors.New("缺少ID参数")
	}

	if err != nil {
		return err
	}

	// 审核成功后广播统计
	if activityId > 0 {
		common.BroadcastCheckInStats(activityId)
	}

	return nil
}

// DeleteCheckIn 删除签到
func (s *AnnualCheckInService) DeleteCheckIn(req annualReq.CheckInDelete) error {
	if len(req.Ids) > 0 {
		return global.GVA_DB.Delete(&annual.AnnualCheckIn{}, "id IN ?", req.Ids).Error
	}
	if req.Id > 0 {
		return global.GVA_DB.Delete(&annual.AnnualCheckIn{}, "id = ?", req.Id).Error
	}
	return errors.New("缺少ID参数")
}

// ExportCheckIn 导出签到
func (s *AnnualCheckInService) ExportCheckIn(info annualReq.CheckInSearch) (filePath string, err error) {
	// TODO: 实现导出逻辑
	return "", nil
}
