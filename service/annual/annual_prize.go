package annual

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualPrizeService struct{}

// CreatePrize 创建奖品
func (s *AnnualPrizeService) CreatePrize(prize annual.AnnualPrize) (err error) {
	return global.GVA_DB.Create(&prize).Error
}

// GetPrizeList 获取奖品列表
func (s *AnnualPrizeService) GetPrizeList(info annualReq.PrizeSearch) (list []annual.AnnualPrize, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualPrize{})

	if info.ActivityId != 0 {
		db = db.Where("activity_id = ?", info.ActivityId)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Level != nil {
		db = db.Where("level = ?", *info.Level)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("level ASC, sort ASC").Find(&list).Error
	return
}

// GetPrizeById 获取奖品详情
func (s *AnnualPrizeService) GetPrizeById(id string) (prize annual.AnnualPrize, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&prize).Error
	return
}

// UpdatePrize 更新奖品
func (s *AnnualPrizeService) UpdatePrize(prize annual.AnnualPrize) (err error) {
	return global.GVA_DB.Model(&annual.AnnualPrize{}).Where("id = ?", prize.ID).Updates(&prize).Error
}

func (s *AnnualPrizeService) DeletePrize(id uint) (err error) {
	if err = s.checkPrizeRelations(id); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&annual.AnnualPrize{}, id).Error
}

// DeletePrizeByIds 批量删除奖品（安全检查）
func (s *AnnualPrizeService) DeletePrizeByIds(ids []uint) (err error) {
	for _, id := range ids {
		if err = s.checkPrizeRelations(id); err != nil {
			return err
		}
	}
	return global.GVA_DB.Delete(&annual.AnnualPrize{}, ids).Error
}

// checkPrizeRelations 检查奖品关联数据
func (s *AnnualPrizeService) checkPrizeRelations(prizeId uint) error {
	var count int64

	// 检查摇一摇场次
	if err := global.GVA_DB.Model(&annual.AnnualShakeRound{}).Where("prize_id = ?", prizeId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该奖品已关联摇一摇场次，无法删除")
	}

	// 检查中奖记录
	if err := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("prize_id = ?", prizeId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该奖品存在中奖记录，无法删除")
	}

	return nil
}
