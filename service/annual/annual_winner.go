package annual

import (
	"errors"
	"go.uber.org/zap"
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
// GetWinnerList 获取中奖列表
func (s *AnnualWinnerService) GetWinnerList(info annualReq.WinnerSearch) (list []annual.AnnualWinner, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualWinner{}).
		Preload("User").
		Preload("Prize").
		Preload("Activity").
		Preload("Round")

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
	if err != nil {
		return
	}

	// ✅ 新增：批量查询摇一摇成绩
	if len(list) > 0 {
		s.fillShakeScores(&list)
	}

	return
}

// fillShakeScores 填充摇一摇成绩信息
func (s *AnnualWinnerService) fillShakeScores(list *[]annual.AnnualWinner) {
	// 收集需要查询的 round_id 和 user_id 组合（只查摇一摇类型）
	type roundUserKey struct {
		RoundId uint
		UserId  uint
	}
	keys := make([]roundUserKey, 0)
	keyIndexMap := make(map[roundUserKey][]int) // 记录每个 key 对应的 list 索引

	for i, winner := range *list {
		// 只有摇一摇类型(winType=1)且有场次ID才查询
		if winner.WinType != nil && *winner.WinType == 1 && winner.RoundId > 0 {
			key := roundUserKey{RoundId: winner.RoundId, UserId: winner.UserId}
			if _, exists := keyIndexMap[key]; !exists {
				keys = append(keys, key)
			}
			keyIndexMap[key] = append(keyIndexMap[key], i)
		}
	}

	if len(keys) == 0 {
		return
	}

	// 构建查询条件
	var scores []annual.AnnualShakeScore
	query := global.GVA_DB.Model(&annual.AnnualShakeScore{})

	// 使用 OR 条件批量查询
	for i, key := range keys {
		if i == 0 {
			query = query.Where("(round_id = ? AND user_id = ?)", key.RoundId, key.UserId)
		} else {
			query = query.Or("(round_id = ? AND user_id = ?)", key.RoundId, key.UserId)
		}
	}

	if err := query.Find(&scores).Error; err != nil {
		global.GVA_LOG.Error("查询摇一摇成绩失败", zap.Error(err))
		return
	}

	// 构建成绩 Map
	scoreMap := make(map[roundUserKey]annual.AnnualShakeScore)
	for _, score := range scores {
		key := roundUserKey{RoundId: score.RoundId, UserId: score.UserId}
		scoreMap[key] = score
	}

	// 回填数据
	for key, indexes := range keyIndexMap {
		if score, exists := scoreMap[key]; exists {
			for _, idx := range indexes {
				(*list)[idx].Score = score.Score
				(*list)[idx].Rank = score.Rank
			}
		}
	}
}

// ConfirmReceive 确认领奖（核销验证）
func (s *AnnualWinnerService) ConfirmReceive(id uint, verifyCode string) (err error) {
	// 1. 查询中奖记录
	var winner annual.AnnualWinner
	if err = global.GVA_DB.Where("id = ?", id).First(&winner).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("中奖记录不存在")
		}
		return err
	}

	// 2. 检查是否已领取
	if winner.Status != nil && *winner.Status == constants.WinnerStatusReceived {
		return errors.New("该奖品已领取，请勿重复核销")
	}

	// 3. 验证核销码
	if winner.ReceiveCode != verifyCode {
		return errors.New("核销码错误，请重新输入")
	}

	// 4. 更新领奖状态
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
