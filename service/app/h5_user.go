package app

import (
	"errors"
	"go.uber.org/zap"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/app/response"
	"lotteryBackend/pkg/cache"
	"lotteryBackend/service/common"
	"time"
)

type H5UserService struct{}

// GetUserInfo 获取用户信息
func (s *H5UserService) GetUserInfo(userId uint, activityId uint) (*response.H5UserResp, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	resp := &response.H5UserResp{
		ID:       user.ID,
		OpenId:   user.OpenId,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		CheckIn:  nil,
	}

	// 如果传了活动ID，查询签到信息
	if activityId > 0 {
		var checkIn annual.AnnualCheckIn
		result := global.GVA_DB.Where("activity_id = ? AND user_id = ?", activityId, userId).First(&checkIn)
		if result.RowsAffected > 0 {
			resp.CheckIn = &response.H5CheckInInfo{
				IsCheckedIn:  true,
				RealName:     checkIn.RealName,
				Phone:        checkIn.Phone,
				Department:   checkIn.Department,
				EmployeeNo:   checkIn.EmployeeNo,
				Status:       checkIn.Status,
				RejectReason: checkIn.RejectReason,
				CheckInTime:  checkIn.CheckInTime,
			}
			// ========== 新增：缓存到 Redis ==========
			go cache.SetUserInfoCache(userId, &cache.UserInfoCache{
				ID:         user.ID,
				Nickname:   user.Nickname,
				Avatar:     user.Avatar,
				RealName:   checkIn.RealName,
				Department: checkIn.Department,
			})
			// ========== 新增结束 ==========
		} else {
			// 未签到，尝试获取上次签到的信息用于自动填充
			var lastCheckIn annual.AnnualCheckIn
			lastResult := global.GVA_DB.Where("user_id = ?", userId).Order("id DESC").First(&lastCheckIn)
			if lastResult.RowsAffected > 0 {
				resp.CheckIn = &response.H5CheckInInfo{
					IsCheckedIn: false,
					RealName:    lastCheckIn.RealName,
					Phone:       lastCheckIn.Phone,
					Department:  lastCheckIn.Department,
					EmployeeNo:  lastCheckIn.EmployeeNo,
					Status:      0,
				}
			} else {
				resp.CheckIn = &response.H5CheckInInfo{
					IsCheckedIn: false,
				}
			}
		}
	}

	return resp, nil
}

// CheckIn 用户签到
func (s *H5UserService) CheckIn(userId uint, req request.CheckInReq, ip string) (*response.H5UserResp, error) {
	// 检查用户是否存在
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查活动是否存在
	var activity annual.AnnualActivity
	if err := global.GVA_DB.First(&activity, req.ActivityId).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	// 检查活动状态
	if *activity.Status != 1 {
		return nil, errors.New("活动未开始或已结束")
	}

	// 检查是否已签到该活动
	var existCheckIn annual.AnnualCheckIn
	result := global.GVA_DB.Where("activity_id = ? AND user_id = ?", req.ActivityId, userId).First(&existCheckIn)
	if result.RowsAffected > 0 {
		return nil, errors.New("您已签到，请勿重复提交")
	}

	// 检查手机号是否在该活动中已被使用
	var phoneCheckIn annual.AnnualCheckIn
	phoneResult := global.GVA_DB.Where("activity_id = ? AND phone = ?", req.ActivityId, req.Phone).First(&phoneCheckIn)
	if phoneResult.RowsAffected > 0 {
		return nil, errors.New("该手机号已被使用")
	}

	// 获取是否需要审核的配置
	needAudit := s.getNeedAuditConfig()
	global.GVA_LOG.Info("是否需要审核", zap.Any("needAudit", needAudit))
	status := 0 // 默认待审核
	// 不需要审核时，直接广播签到统计
	if !needAudit {
		status = 1
		common.BroadcastCheckInStats(req.ActivityId)
	}

	// 创建签到记录
	checkIn := annual.AnnualCheckIn{
		ActivityId:  req.ActivityId,
		UserId:      userId,
		RealName:    req.RealName,
		Phone:       req.Phone,
		Department:  req.Department,
		EmployeeNo:  req.EmployeeNo,
		Status:      status,
		CheckInTime: time.Now(),
		Ip:          ip,
	}

	if err := global.GVA_DB.Create(&checkIn).Error; err != nil {
		return nil, errors.New("签到失败，请稍后重试")
	}

	return &response.H5UserResp{
		ID:       user.ID,
		OpenId:   user.OpenId,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		CheckIn: &response.H5CheckInInfo{
			IsCheckedIn:  true,
			RealName:     checkIn.RealName,
			Phone:        checkIn.Phone,
			Department:   checkIn.Department,
			EmployeeNo:   checkIn.EmployeeNo,
			Status:       checkIn.Status,
			RejectReason: "",
			CheckInTime:  checkIn.CheckInTime,
		},
	}, nil
}

// getNeedAuditConfig 获取是否需要审核配置
func (s *H5UserService) getNeedAuditConfig() bool {
	var config annual.AnnualConfig
	global.GVA_DB.Where("config_key = ?", "need_audit").First(&config)
	return config.ConfigValue == "1" || config.ConfigValue == "true"
}

// GetUserBrief 获取用户简要信息
func (s *H5UserService) GetUserBrief(userId uint) *response.H5UserBrief {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil
	}

	return &response.H5UserBrief{
		ID:       user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}
}

// GetCheckInBrief 获取签到用户简要信息（用于列表）
func (s *H5UserService) GetCheckInBrief(checkIn *annual.AnnualCheckIn) *response.H5UserBrief {
	var user annual.AnnualUser
	global.GVA_DB.First(&user, checkIn.UserId)

	return &response.H5UserBrief{
		ID:         user.ID,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		RealName:   checkIn.RealName,
		Department: checkIn.Department,
	}
}

// CheckUserCanJoin 检查用户是否可以参与活动（已签到且审核通过）
func (s *H5UserService) CheckUserCanJoin(userId uint, activityId uint) error {
	var checkIn annual.AnnualCheckIn
	result := global.GVA_DB.Where("activity_id = ? AND user_id = ?", activityId, userId).First(&checkIn)

	if result.RowsAffected == 0 {
		return errors.New("请先完成签到")
	}

	if checkIn.Status == 0 {
		return errors.New("您的签到正在审核中")
	}

	if checkIn.Status == 2 {
		return errors.New("您的签到未通过审核")
	}

	return nil
}
