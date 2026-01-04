package app

import (
	"encoding/json"
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
)

type H5PrizeService struct{}

// GetPrizeList 获取奖品列表
func (s *H5PrizeService) GetPrizeList(activityId uint) ([]response.H5PrizeResp, error) {
	var prizes []annual.AnnualPrize
	err := global.GVA_DB.Where("activity_id = ?", activityId).
		Order("level ASC, sort ASC").
		Find(&prizes).Error

	if err != nil {
		return nil, err
	}

	result := make([]response.H5PrizeResp, len(prizes))
	for i, p := range prizes {
		result[i] = response.H5PrizeResp{
			ID:          p.ID,
			Name:        p.Name,
			Image:       p.Image,
			Level:       *p.Level,
			TotalCount:  p.TotalCount,
			RemainCount: p.RemainCount,
		}
	}

	return result, nil
}

// GetPrizeDetail 获取奖品详情
func (s *H5PrizeService) GetPrizeDetail(prizeId uint) (*response.H5PrizeResp, error) {
	var prize annual.AnnualPrize
	if err := global.GVA_DB.First(&prize, prizeId).Error; err != nil {
		return nil, errors.New("奖品不存在")
	}

	return &response.H5PrizeResp{
		ID:          prize.ID,
		Name:        prize.Name,
		Image:       prize.Image,
		Level:       *prize.Level,
		TotalCount:  prize.TotalCount,
		RemainCount: prize.RemainCount,
	}, nil
}

// GetMyWinnings 获取我的中奖记录
func (s *H5PrizeService) GetMyWinnings(userId, activityId uint) ([]response.WinningResp, error) {
	var winners []annual.AnnualWinner
	err := global.GVA_DB.Where("activity_id = ? AND user_id = ?", activityId, userId).
		Order("created_at DESC").
		Find(&winners).Error

	if err != nil {
		return nil, err
	}

	result := make([]response.WinningResp, len(winners))
	for i, w := range winners {
		result[i] = response.WinningResp{
			ID:          w.ID,
			PrizeId:     w.PrizeId,
			WinType:     *w.WinType,
			Status:      *w.Status,
			ReceiveTime: w.ReceiveTime,
			CreatedAt:   w.CreatedAt,
		}

		// 获取奖品信息
		result[i].Prize = s.getPrizeBrief(w.PrizeId)
	}

	return result, nil
}

// GetWinningDetail 获取中奖详情
func (s *H5PrizeService) GetWinningDetail(winnerId uint) (*response.WinningResp, error) {
	var winner annual.AnnualWinner
	if err := global.GVA_DB.First(&winner, winnerId).Error; err != nil {
		return nil, errors.New("中奖记录不存在")
	}

	userService := H5UserService{}
	userBrief := userService.GetUserBrief(winner.UserId)

	resp := &response.WinningResp{
		ID:          winner.ID,
		PrizeId:     winner.PrizeId,
		WinType:     *winner.WinType,
		Status:      *winner.Status,
		ReceiveTime: winner.ReceiveTime,
		CreatedAt:   winner.CreatedAt,
		Prize:       s.getPrizeBrief(winner.PrizeId),
	}

	if userBrief != nil {
		resp.User = userBrief
	}

	return resp, nil
}

// GetRecentWinnings 获取最新中奖记录
func (s *H5PrizeService) GetRecentWinnings(activityId uint, limit int) ([]response.WinningResp, error) {
	if limit <= 0 {
		limit = 10
	}

	var winners []annual.AnnualWinner
	err := global.GVA_DB.Where("activity_id = ?", activityId).
		Order("created_at DESC").
		Limit(limit).
		Find(&winners).Error

	if err != nil {
		return nil, err
	}

	userService := H5UserService{}
	result := make([]response.WinningResp, len(winners))
	for i, w := range winners {
		userBrief := userService.GetUserBrief(w.UserId)

		result[i] = response.WinningResp{
			ID:          w.ID,
			PrizeId:     w.PrizeId,
			WinType:     *w.WinType,
			Status:      *w.Status,
			ReceiveTime: w.ReceiveTime,
			CreatedAt:   w.CreatedAt,
			Prize:       s.getPrizeBrief(w.PrizeId),
		}

		if userBrief != nil {
			result[i].User = userBrief
		}
	}

	return result, nil
}

// GetReceiveQrCode 获取领奖二维码数据
func (s *H5PrizeService) GetReceiveQrCode(userId, winnerId uint) (string, error) {
	var winner annual.AnnualWinner
	if err := global.GVA_DB.First(&winner, winnerId).Error; err != nil {
		return "", errors.New("中奖记录不存在")
	}

	// 验证是否是本人的中奖记录
	if winner.UserId != userId {
		return "", errors.New("无权查看")
	}

	// 检查是否已领取
	if *winner.Status == 1 {
		return "", errors.New("奖品已领取")
	}

	// 生成二维码数据
	qrData := map[string]interface{}{
		"type":     "prize_receive",
		"winnerId": winner.ID,
		"prizeId":  winner.PrizeId,
		"userId":   winner.UserId,
	}

	qrBytes, _ := json.Marshal(qrData)
	return string(qrBytes), nil
}

// getPrizeBrief 获取奖品简要信息
func (s *H5PrizeService) getPrizeBrief(prizeId uint) *response.H5PrizeBrief {
	if prizeId == 0 {
		return nil
	}

	var prize annual.AnnualPrize
	if err := global.GVA_DB.First(&prize, prizeId).Error; err != nil {
		return nil
	}

	return &response.H5PrizeBrief{
		ID:    prize.ID,
		Name:  prize.Name,
		Image: prize.Image,
		Level: *prize.Level,
	}
}

// ConfirmReceive 确认领奖（供后台调用）
func (s *H5PrizeService) ConfirmReceive(winnerId uint) error {
	var winner annual.AnnualWinner
	if err := global.GVA_DB.First(&winner, winnerId).Error; err != nil {
		return errors.New("中奖记录不存在")
	}

	if *winner.Status == 1 {
		return errors.New("奖品已领取")
	}

	return global.GVA_DB.Model(&winner).Updates(map[string]interface{}{
		"status":       1,
		"receive_time": global.GVA_DB.Raw("NOW()"),
	}).Error
}
