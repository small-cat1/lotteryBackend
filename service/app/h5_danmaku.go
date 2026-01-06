package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/response"
)

type H5DanmakuService struct{}

// SendDanmaku 发送弹幕
func (s *H5DanmakuService) SendDanmaku(userId, activityId uint, content, color string) (*response.DanmakuResp, error) {
	// 检查活动
	var activity annual.AnnualActivity
	if err := global.GVA_DB.First(&activity, activityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	if *activity.Status != 1 {
		return nil, errors.New("活动未开始或已结束")
	}

	if *activity.DanmakuEnabled != 1 {
		return nil, errors.New("弹幕功能未开启")
	}

	// 内容长度检查
	if len(content) == 0 {
		return nil, errors.New("弹幕内容不能为空")
	}
	if len(content) > 100 {
		return nil, errors.New("弹幕内容过长")
	}

	// 默认颜色
	if color == "" {
		color = "#FFFFFF"
	}

	// 设置审核状态
	status := 1 // 默认通过
	if *activity.DanmakuAudit == 1 {
		status = 0 // 待审核
	}

	// 创建弹幕
	danmaku := annual.AnnualDanmaku{
		ActivityId: activityId,
		UserId:     userId,
		Content:    content,
		Color:      color,
		Status:     &status,
	}

	if err := global.GVA_DB.Create(&danmaku).Error; err != nil {
		return nil, errors.New("发送失败，请稍后重试")
	}

	// 获取用户信息
	userService := H5UserService{}
	userBrief := userService.GetUserBrief(userId)

	resp := &response.DanmakuResp{
		ID:        danmaku.ID,
		Content:   danmaku.Content,
		Color:     danmaku.Color,
		IsTop:     *danmaku.IsTop,
		CreatedAt: danmaku.CreatedAt,
	}
	if userBrief != nil {
		resp.User = *userBrief
	}

	return resp, nil
}

// GetDanmakuList 获取弹幕列表
func (s *H5DanmakuService) GetDanmakuList(activityId uint, page, pageSize int) (*response.H5PageResult, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var danmakus []annual.AnnualDanmaku
	var total int64

	db := global.GVA_DB.Model(&annual.AnnualDanmaku{}).
		Where("activity_id = ? AND status = ?", activityId, 1)
	db.Count(&total)

	err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&danmakus).Error

	if err != nil {
		return nil, err
	}

	userService := H5UserService{}
	list := make([]response.DanmakuResp, len(danmakus))
	for i, d := range danmakus {
		userBrief := userService.GetUserBrief(d.UserId)
		list[i] = response.DanmakuResp{
			ID:        d.ID,
			Content:   d.Content,
			Color:     d.Color,
			IsTop:     *d.IsTop,
			CreatedAt: d.CreatedAt,
		}
		if userBrief != nil {
			list[i].User = *userBrief
		}
	}

	return &response.H5PageResult{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
