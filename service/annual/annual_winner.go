package annual

import (
	"errors"
	"lotteryBackend/constants"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"math/rand"
	"time"

	"gorm.io/gorm"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualWinnerService struct{}

// GetWinnerList 获取中奖列表
func (s *AnnualWinnerService) GetWinnerList(info annualReq.WinnerSearch) (list []annual.AnnualWinner, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualWinner{}).
		Preload("User").
		Preload("Prize").
		Preload("Activity")

	if info.ActivityId != 0 {
		db = db.Where("activity_id = ?", info.ActivityId)
	}
	if info.UserId != 0 {
		db = db.Where("user_id = ?", info.UserId)
	}
	if info.PrizeId != 0 {
		db = db.Where("prize_id = ?", info.PrizeId)
	}
	if info.RoundId != 0 {
		db = db.Where("round_id = ?", info.RoundId)
	}
	if info.WinType != nil {
		db = db.Where("win_type = ?", *info.WinType)
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

// ConfirmReceive 确认领奖
func (s *AnnualWinnerService) ConfirmReceive(id uint) (err error) {
	now := time.Now()
	return global.GVA_DB.Model(&annual.AnnualWinner{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       constants.WinnerStatusReceived,
			"receive_time": now,
		}).Error
}

// DeleteWinner 删除中奖记录
func (s *AnnualWinnerService) DeleteWinner(id uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualWinner{}, id).Error
}

// ExportWinner 导出中奖记录
func (s *AnnualWinnerService) ExportWinner(info annualReq.WinnerSearch) (filePath string, err error) {
	// TODO: 实现导出逻辑
	return "", nil
}

// RandomDraw 随机抽奖
func (s *AnnualWinnerService) RandomDraw(req annualReq.RandomDraw) (winners []annual.AnnualUser, err error) {
	// 获取活动配置
	var activity annual.AnnualActivity
	if err = global.GVA_DB.First(&activity, req.ActivityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	// 检查奖品库存
	var prize annual.AnnualPrize
	if err = global.GVA_DB.First(&prize, req.PrizeId).Error; err != nil {
		return nil, errors.New("奖品不存在")
	}
	if prize.RemainCount < req.Count {
		return nil, errors.New("奖品库存不足")
	}

	// 构建候选人查询
	db := global.GVA_DB.Model(&annual.AnnualUser{}).
		Joins("INNER JOIN annual_check_ins ON annual_users.id = annual_check_ins.user_id").
		Where("annual_check_ins.activity_id = ?", req.ActivityId).
		Where("annual_users.is_registered = ?", constants.UserRegistered).
		Where("annual_users.status = ?", constants.UserStatusApproved)

	// 如果配置了中奖后排除
	if *activity.WinnerExclude == 1 {
		db = db.Where("annual_users.id NOT IN (?)",
			global.GVA_DB.Model(&annual.AnnualWinner{}).
				Select("user_id").
				Where("activity_id = ?", req.ActivityId))
	}

	// 获取候选人
	var candidates []annual.AnnualUser
	if err = db.Find(&candidates).Error; err != nil {
		return
	}

	if len(candidates) < req.Count {
		return nil, errors.New("候选人数量不足")
	}

	// 随机抽取
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	winners = candidates[:req.Count]

	// 开启事务保存中奖记录
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, winner := range winners {
			winType := constants.WinTypeRandom
			record := annual.AnnualWinner{
				ActivityId: req.ActivityId,
				UserId:     winner.ID,
				PrizeId:    req.PrizeId,
				WinType:    &winType,
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}

		// 扣减库存
		if err := tx.Model(&annual.AnnualPrize{}).
			Where("id = ? AND remain_count >= ?", req.PrizeId, req.Count).
			Update("remain_count", gorm.Expr("remain_count - ?", req.Count)).Error; err != nil {
			return err
		}

		return nil
	})

	return
}
