package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualDanmakuService struct{}

// GetDanmakuList 获取弹幕列表
func (s *AnnualDanmakuService) GetDanmakuList(info annualReq.DanmakuSearch) (list []annual.AnnualDanmaku, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualDanmaku{}).Preload("User")

	if info.ActivityId != 0 {
		db = db.Where("activity_id = ?", info.ActivityId)
	}
	if info.UserId != 0 {
		db = db.Where("user_id = ?", info.UserId)
	}
	if info.Content != "" {
		db = db.Where("content LIKE ?", "%"+info.Content+"%")
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

// AuditDanmaku 审核弹幕
func (s *AnnualDanmakuService) AuditDanmaku(ids []uint, status int) (err error) {
	return global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("id IN ?", ids).Update("status", status).Error
}

// TopDanmaku 置顶弹幕
func (s *AnnualDanmakuService) TopDanmaku(id uint, isTop int) (err error) {
	return global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("id = ?", id).Update("is_top", isTop).Error
}

// DeleteDanmaku 删除弹幕
func (s *AnnualDanmakuService) DeleteDanmaku(id uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualDanmaku{}, id).Error
}

// DeleteDanmakuByIds 批量删除弹幕
func (s *AnnualDanmakuService) DeleteDanmakuByIds(ids []uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualDanmaku{}, ids).Error
}
