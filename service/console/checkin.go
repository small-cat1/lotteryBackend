package console

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/console/response"
	"lotteryBackend/service/common"
)

type CheckInService struct{}

// GetCheckInStats 获取签到统计（包含状态）
func (s *CheckInService) GetCheckInStats(activityId uint) (*response.CheckInStatsResp, error) {
	// 获取签到开关状态
	isOpen, err := common.GetCheckInSwitch(activityId)
	if err != nil {
		isOpen = false
	}

	// 从Redis获取已签到人数
	checkedCount, err := common.GetCheckInCount(activityId)
	if err != nil {
		// Redis出错，降级查数据库
		global.GVA_DB.Model(&annual.AnnualCheckIn{}).
			Where("activity_id = ?", activityId).
			Count(&checkedCount)
	}

	// 总人数（已报名且通过审核）
	var totalCount int64
	global.GVA_DB.Model(&annual.AnnualUser{}).
		Where("status = ?", 1).
		Count(&totalCount)

	var rate float64 = 0
	if totalCount > 0 {
		rate = float64(checkedCount) / float64(totalCount) * 100
	}

	return &response.CheckInStatsResp{
		IsOpen:    isOpen,
		CheckedIn: int(checkedCount),
		Total:     int(totalCount),
		Rate:      rate,
	}, nil
}

// OpenCheckIn 开启签到
func (s *CheckInService) OpenCheckIn(activityId uint) error {
	// 检查活动是否存在（从缓存或数据库）
	activityService := ActivityService{}
	_, err := activityService.GetActivityDetail(activityId)
	if err != nil {
		return err
	}

	// 设置签到开关为开启
	return common.SetCheckInSwitch(activityId, true)
}

// CloseCheckIn 关闭签到
func (s *CheckInService) CloseCheckIn(activityId uint) error {
	return common.SetCheckInSwitch(activityId, false)
}

// GetCheckInList 获取签到列表
func (s *CheckInService) GetCheckInList(activityId uint, page, pageSize int, keyword string) (*response.CheckInListResp, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	query := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("activity_id = ?", activityId)

	// 关键词搜索
	if keyword != "" {
		var userIds []uint
		global.GVA_DB.Model(&annual.AnnualUser{}).
			Where("nickname LIKE ? OR real_name LIKE ? OR phone LIKE ?",
				"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
			Pluck("id", &userIds)
		if len(userIds) > 0 {
			query = query.Where("user_id IN ?", userIds)
		} else {
			return &response.CheckInListResp{
				List:     []response.CheckInItem{},
				Total:    0,
				Page:     page,
				PageSize: pageSize,
			}, nil
		}
	}

	var total int64
	query.Count(&total)

	var checkIns []annual.AnnualCheckIn
	query.Order("check_in_time DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&checkIns)

	list := make([]response.CheckInItem, 0, len(checkIns))
	for _, c := range checkIns {
		item := response.CheckInItem{
			ID:          c.ID,
			UserId:      c.UserId,
			CheckInTime: c.CheckInTime,
		}

		// 获取用户信息
		var user annual.AnnualUser
		if err := global.GVA_DB.First(&user, c.UserId).Error; err == nil {
			item.User = &response.UserInfo{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   user.RealName,
				Department: user.Department,
				Phone:      user.Phone,
			}
		}

		list = append(list, item)
	}

	return &response.CheckInListResp{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// SyncCheckInUsersToRedis 同步签到用户到Redis（启动时或数据恢复用）
func (s *CheckInService) SyncCheckInUsersToRedis(activityId uint) error {
	var userIds []uint
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ?", activityId).
		Pluck("user_id", &userIds)

	for _, userId := range userIds {
		common.AddCheckInUser(activityId, userId)
	}

	return nil
}

// IsUserCheckedIn 检查用户是否已签到
func (s *CheckInService) IsUserCheckedIn(activityId uint, userId uint) (bool, error) {
	// 先查Redis
	checked, err := common.IsUserCheckIn(activityId, userId)
	if err == nil {
		return checked, nil
	}

	// Redis出错，降级查数据库
	var count int64
	global.GVA_DB.Model(&annual.AnnualCheckIn{}).
		Where("activity_id = ? AND user_id = ?", activityId, userId).
		Count(&count)
	return count > 0, nil
}

// DoCheckIn 执行签到（供H5端调用）
func (s *CheckInService) DoCheckIn(activityId uint, userId uint, ip string) error {
	// 检查签到是否开启
	isOpen, _ := common.GetCheckInSwitch(activityId)
	if !isOpen {
		return errors.New("签到未开启")
	}

	// 检查是否已签到
	checked, _ := s.IsUserCheckedIn(activityId, userId)
	if checked {
		return errors.New("您已签到")
	}

	// 添加到Redis
	added, err := common.AddCheckInUser(activityId, userId)
	if err != nil {
		return errors.New("签到失败")
	}
	if !added {
		return errors.New("您已签到")
	}

	// 写入数据库
	checkIn := annual.AnnualCheckIn{
		ActivityId: activityId,
		UserId:     userId,
		Ip:         ip,
	}
	if err := global.GVA_DB.Create(&checkIn).Error; err != nil {
		return errors.New("签到失败")
	}

	return nil
}
