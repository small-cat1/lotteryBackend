package annual

import (
	"errors"
	"lotteryBackend/constants"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"strconv"
	"time"

	"gorm.io/gorm"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualShakeRoundService struct{}

// CreateShakeRound 创建摇一摇场次
func (s *AnnualShakeRoundService) CreateShakeRound(round annual.AnnualShakeRound) (err error) {
	return global.GVA_DB.Create(&round).Error
}

// GetShakeRoundList 获取场次列表
func (s *AnnualShakeRoundService) GetShakeRoundList(info annualReq.ShakeRoundSearch) (list []annual.AnnualShakeRound, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualShakeRound{}).Preload("Prize")

	if info.ActivityId != 0 {
		db = db.Where("activity_id = ?", info.ActivityId)
	}
	if info.RoundName != "" {
		db = db.Where("round_name LIKE ?", "%"+info.RoundName+"%")
	}
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("sort ASC, id DESC").Find(&list).Error
	return
}

// GetShakeRoundById 获取场次详情
func (s *AnnualShakeRoundService) GetShakeRoundById(id string) (round annual.AnnualShakeRound, err error) {
	err = global.GVA_DB.Preload("Prize").Where("id = ?", id).First(&round).Error
	return
}

// UpdateShakeRound 更新场次
func (s *AnnualShakeRoundService) UpdateShakeRound(round annual.AnnualShakeRound) (err error) {
	return global.GVA_DB.Model(&round).
		Select("activity_id", "round_name", "duration", "winner_count", "prize_id", "sort").
		Updates(&round).Error
}

// StartShakeRound 开始游戏
func (s *AnnualShakeRoundService) StartShakeRound(id string) (err error) {
	var round annual.AnnualShakeRound
	if err = global.GVA_DB.First(&round, id).Error; err != nil {
		return
	}
	if *round.Status != constants.ShakeRoundStatusNotStarted {
		return errors.New("游戏状态不正确")
	}

	now := time.Now()
	status := constants.ShakeRoundStatusOngoing
	return global.GVA_DB.Model(&round).Updates(map[string]interface{}{
		"status":     status,
		"start_time": now,
	}).Error
}

// StopShakeRound 结束游戏并计算获奖者
func (s *AnnualShakeRoundService) StopShakeRound(id string) (winners []annual.AnnualShakeScore, err error) {
	var round annual.AnnualShakeRound
	if err = global.GVA_DB.Preload("Prize").First(&round, id).Error; err != nil {
		return
	}
	if *round.Status != constants.ShakeRoundStatusOngoing {
		return nil, errors.New("游戏状态不正确")
	}

	// 开启事务
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		status := constants.ShakeRoundStatusEnded

		// 更新场次状态
		if err := tx.Model(&round).Updates(map[string]interface{}{
			"status":   status,
			"end_time": now,
		}).Error; err != nil {
			return err
		}

		// 获取排行榜前N名
		if err := tx.Where("round_id = ?", id).
			Order("score DESC").
			Limit(round.WinnerCount).
			Preload("User").
			Find(&winners).Error; err != nil {
			return err
		}

		// 更新排名和获奖状态
		for i, winner := range winners {
			rank := i + 1
			isWinner := 1
			if err := tx.Model(&annual.AnnualShakeScore{}).
				Where("id = ?", winner.ID).
				Updates(map[string]interface{}{
					"rank":      rank,
					"is_winner": isWinner,
				}).Error; err != nil {
				return err
			}
			winners[i].Rank = rank
			winners[i].IsWinner = &isWinner

			// 创建中奖记录
			roundIdUint, _ := strconv.ParseUint(id, 10, 64)
			winType := constants.WinTypeShake
			winnerRecord := annual.AnnualWinner{
				ActivityId: round.ActivityId,
				UserId:     winner.UserId,
				PrizeId:    round.PrizeId,
				RoundId:    uint(roundIdUint),
				WinType:    &winType,
			}
			if err := tx.Create(&winnerRecord).Error; err != nil {
				return err
			}
		}

		// 扣减奖品库存
		if round.PrizeId > 0 && len(winners) > 0 {
			if err := tx.Model(&annual.AnnualPrize{}).
				Where("id = ? AND remain_count >= ?", round.PrizeId, len(winners)).
				Update("remain_count", gorm.Expr("remain_count - ?", len(winners))).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return
}

// DeleteShakeRound 删除摇一摇场次（安全检查）
func (s *AnnualShakeRoundService) DeleteShakeRound(id uint) (err error) {
	if err = s.checkShakeRoundRelations(id); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&annual.AnnualShakeRound{}, id).Error
}

// checkShakeRoundRelations 检查场次关联数据
func (s *AnnualShakeRoundService) checkShakeRoundRelations(roundId uint) error {
	var count int64

	// 检查摇一摇成绩
	if err := global.GVA_DB.Model(&annual.AnnualShakeScore{}).Where("round_id = ?", roundId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该场次存在成绩记录，无法删除")
	}

	// 检查中奖记录
	if err := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("round_id = ?", roundId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该场次存在中奖记录，无法删除")
	}

	return nil
}

// GetShakeScores 获取成绩排行
func (s *AnnualShakeRoundService) GetShakeScores(roundId string) (scores []annual.AnnualShakeScore, err error) {
	err = global.GVA_DB.Where("round_id = ?", roundId).
		Order("score DESC").
		Preload("User").
		Find(&scores).Error
	return
}
