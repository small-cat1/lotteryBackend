package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/service/common"
)

type DanmakuService struct{}

// GetDanmakuList 获取弹幕列表
func (s *DanmakuService) GetDanmakuList(activityId uint, limit int, status int) (*response.DanmakuListResp, error) {
	if limit <= 0 {
		limit = 50
	}

	query := global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("activity_id = ?", activityId)
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	var danmakus []annual.AnnualDanmaku
	query.Order("created_at DESC").Limit(limit).Find(&danmakus)

	// 获取弹幕用户数（去重）
	var userCount int64
	countQuery := global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("activity_id = ?", activityId)
	if status > 0 {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Distinct("user_id").Count(&userCount)

	list := make([]response.DanmakuItem, 0, len(danmakus))
	for _, d := range danmakus {
		item := response.DanmakuItem{
			ID:        d.ID,
			Content:   d.Content,
			Color:     d.Color,
			Status:    safeInt(d.Status),
			IsTop:     safeInt(d.IsTop),
			CreatedAt: d.CreatedAt,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, d.UserId).Error; err == nil {
			item.User = &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		list = append(list, item)
	}

	return &response.DanmakuListResp{
		List:      list,
		UserCount: int(userCount),
	}, nil
}

// AuditDanmaku 审核弹幕
func (s *DanmakuService) AuditDanmaku(danmakuId uint, status int) error {
	var danmaku annual.AnnualDanmaku
	if err := global.GVA_DB.First(&danmaku, danmakuId).Error; err != nil {
		return errors.New("弹幕不存在")
	}

	return global.GVA_DB.Model(&danmaku).Update("status", status).Error
}

// OpenDanmaku 开启弹幕
func (s *DanmakuService) OpenDanmaku(activityId uint) error {
	return common.SetDanmakuSwitch(activityId, true)
}

// CloseDanmaku 关闭弹幕
func (s *DanmakuService) CloseDanmaku(activityId uint) error {
	return common.SetDanmakuSwitch(activityId, false)
}

// IsDanmakuOpen 检查弹幕是否开启
func (s *DanmakuService) IsDanmakuOpen(activityId uint) (bool, error) {
	return common.GetDanmakuSwitch(activityId)
}

// GetDanmakuUsers 获取发送弹幕的用户列表（用于弹幕抽奖）
func (s *DanmakuService) GetDanmakuUsers(activityId uint) ([]uint, error) {
	var userIds []uint
	global.GVA_DB.Model(&annual.AnnualDanmaku{}).
		Where("activity_id = ? AND status = ?", activityId, 1).
		Distinct().
		Pluck("user_id", &userIds)
	return userIds, nil
}
