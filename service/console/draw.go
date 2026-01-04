package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/service/common"
	"math/rand"
	"time"
)

type DrawService struct{}

// 中奖类型
const (
	WinTypeShake   = 1 // 摇一摇
	WinTypeRandom  = 2 // 随机抽奖
	WinTypeDanmaku = 3 // 弹幕抽奖
)

// RandomDraw 随机抽奖（从已签到用户中抽取）
func (s *DrawService) RandomDraw(activityId uint, count int, prizeId uint) (*response.DrawResultResp, error) {
	// 获取已签到用户
	userIdStrs, err := common.GetCheckInUserIds(activityId)
	if err != nil || len(userIdStrs) == 0 {
		// 降级查数据库
		var userIds []uint
		global.GVA_DB.Model(&annual.AnnualCheckIn{}).
			Where("activity_id = ?", activityId).
			Pluck("user_id", &userIds)
		if len(userIds) == 0 {
			return nil, errors.New("暂无签到用户")
		}
		return s.doDraw(activityId, userIds, count, prizeId, WinTypeRandom)
	}

	// 转换用户ID
	userIds := make([]uint, 0, len(userIdStrs))
	for _, idStr := range userIdStrs {
		var id uint
		if _, err := parseUint(idStr, &id); err == nil {
			userIds = append(userIds, id)
		}
	}

	if len(userIds) == 0 {
		return nil, errors.New("暂无签到用户")
	}

	return s.doDraw(activityId, userIds, count, prizeId, WinTypeRandom)
}

// DanmakuDraw 弹幕抽奖
func (s *DrawService) DanmakuDraw(activityId uint, count int, prizeId uint) (*response.DrawResultResp, error) {
	danmakuService := DanmakuService{}
	userIds, err := danmakuService.GetDanmakuUsers(activityId)
	if err != nil || len(userIds) == 0 {
		return nil, errors.New("暂无弹幕用户")
	}

	return s.doDraw(activityId, userIds, count, prizeId, WinTypeDanmaku)
}

// doDraw 执行抽奖
func (s *DrawService) doDraw(activityId uint, candidates []uint, count int, prizeId uint, winType int) (*response.DrawResultResp, error) {
	// 检查活动是否排除已中奖用户
	activityService := ActivityService{}
	activityDetail, _ := activityService.GetActivityDetail(activityId)
	excludeWinner := activityDetail != nil && activityDetail.WinnerExclude == 1

	// 过滤已中奖用户
	if excludeWinner {
		candidates = s.filterWonUsers(activityId, candidates)
	}

	if len(candidates) < count {
		return nil, errors.New("符合条件的用户不足")
	}

	// 获取奖品
	if prizeId == 0 {
		return nil, errors.New("未指定奖品")
	}

	var prize annual.AnnualPrize
	if err := global.GVA_DB.First(&prize, prizeId).Error; err != nil {
		return nil, errors.New("奖品不存在")
	}

	if prize.RemainCount < count {
		return nil, errors.New("奖品库存不足")
	}

	// 随机抽取
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	selectedIds := candidates[:count]

	// 创建中奖记录
	prizeInfo := &response.PrizeItem{
		ID:          prize.ID,
		Name:        prize.Name,
		Image:       prize.Image,
		Level:       safeInt(prize.Level),
		TotalCount:  prize.TotalCount,
		RemainCount: prize.RemainCount - count,
	}

	winners := make([]response.WinnerItem, 0, count)

	for _, userId := range selectedIds {
		var user annual.AnnualUser
		global.GVA_DB.First(&user, userId)

		winner := annual.AnnualWinner{
			ActivityId: activityId,
			UserId:     userId,
			PrizeId:    prizeId,
			WinType:    &winType,
		}
		global.GVA_DB.Create(&winner)

		userInfo := &response.UserInfo{
			ID:         user.ID,
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			RealName:   user.RealName,
			Department: user.Department,
		}

		winners = append(winners, response.WinnerItem{
			ID:        winner.ID,
			UserId:    userId,
			WinType:   winType,
			CreatedAt: winner.CreatedAt,
			User:      userInfo,
			Prize:     prizeInfo,
		})
	}

	// 更新奖品剩余数量
	global.GVA_DB.Model(&prize).Update("remain_count", prize.RemainCount-count)

	return &response.DrawResultResp{
		Winners: winners,
		Prize:   prizeInfo,
	}, nil
}

// filterWonUsers 过滤已中奖用户
func (s *DrawService) filterWonUsers(activityId uint, userIds []uint) []uint {
	var wonUserIds []uint
	global.GVA_DB.Model(&annual.AnnualWinner{}).
		Where("activity_id = ?", activityId).
		Pluck("user_id", &wonUserIds)

	wonMap := make(map[uint]bool)
	for _, id := range wonUserIds {
		wonMap[id] = true
	}

	candidates := make([]uint, 0)
	for _, id := range userIds {
		if !wonMap[id] {
			candidates = append(candidates, id)
		}
	}

	return candidates
}

// GetAllWinners 获取活动所有中奖名单
func (s *DrawService) GetAllWinners(activityId uint, winType int) (*response.WinnerListResp, error) {
	query := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("activity_id = ?", activityId)
	if winType > 0 {
		query = query.Where("win_type = ?", winType)
	}

	var winners []annual.AnnualWinner
	query.Order("created_at DESC").Find(&winners)

	list := make([]response.WinnerItem, 0, len(winners))
	for _, w := range winners {
		item := response.WinnerItem{
			ID:        w.ID,
			UserId:    w.UserId,
			WinType:   safeInt(w.WinType),
			CreatedAt: w.CreatedAt,
		}

		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, w.UserId).Error; err == nil {
			item.User = &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
			}
		}

		var prize annual.AnnualPrize
		if err := global.GVA_DB.First(&prize, w.PrizeId).Error; err == nil {
			item.Prize = &response.PrizeItem{
				ID:    prize.ID,
				Name:  prize.Name,
				Image: prize.Image,
				Level: safeInt(prize.Level),
			}
		}

		list = append(list, item)
	}

	return &response.WinnerListResp{List: list}, nil
}

// parseUint 解析uint
func parseUint(s string, result *uint) (int, error) {
	var val uint
	n, err := sscanf(s, "%d", &val)
	if err == nil {
		*result = val
	}
	return n, err
}

// sscanf 简单的sscanf实现
func sscanf(s string, format string, args ...interface{}) (int, error) {
	var val uint
	for _, c := range s {
		if c >= '0' && c <= '9' {
			val = val*10 + uint(c-'0')
		}
	}
	if len(args) > 0 {
		if p, ok := args[0].(*uint); ok {
			*p = val
			return 1, nil
		}
	}
	return 0, nil
}
