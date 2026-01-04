package app

import (
	"errors"
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	"lotteryBackend/model/app/request"
	"lotteryBackend/model/app/response"
)

type H5UserService struct{}

// GetUserInfo 获取用户信息
func (s *H5UserService) GetUserInfo(userId uint) (*response.H5UserResp, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	isRegistered := 0
	if user.RealName != "" {
		isRegistered = 1
	}

	return &response.H5UserResp{
		ID:           user.ID,
		OpenId:       user.OpenId,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		RealName:     user.RealName,
		Phone:        user.Phone,
		Department:   user.Department,
		EmployeeNo:   user.EmployeeNo,
		IsRegistered: isRegistered,
		Status:       *user.Status,
	}, nil
}

// UserRegister 用户报名
func (s *H5UserService) UserRegister(userId uint, req request.UserRegisterReq) (*response.H5UserResp, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查是否已报名
	if user.RealName != "" {
		return nil, errors.New("您已报名，请勿重复提交")
	}

	// 检查手机号是否已被使用
	var existUser annual.AnnualUser
	result := global.GVA_DB.Where("phone = ? AND id != ?", req.Phone, userId).First(&existUser)
	if result.RowsAffected > 0 {
		return nil, errors.New("该手机号已被使用")
	}

	// 检查工号是否已被使用
	result = global.GVA_DB.Where("employee_no = ? AND id != ?", req.EmployeeNo, userId).First(&existUser)
	if result.RowsAffected > 0 {
		return nil, errors.New("该工号已被使用")
	}

	// 获取是否需要审核的配置
	needAudit := s.getNeedAuditConfig()

	// 更新用户信息
	status := 1 // 默认通过
	if needAudit {
		status = 0 // 待审核
	}

	err := global.GVA_DB.Model(&user).Updates(map[string]interface{}{
		"real_name":   req.RealName,
		"phone":       req.Phone,
		"department":  req.Department,
		"employee_no": req.EmployeeNo,
		"status":      status,
	}).Error

	if err != nil {
		return nil, errors.New("报名失败，请稍后重试")
	}

	// 重新查询用户信息
	global.GVA_DB.First(&user, userId)

	return &response.H5UserResp{
		ID:           user.ID,
		OpenId:       user.OpenId,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		RealName:     user.RealName,
		Phone:        user.Phone,
		Department:   user.Department,
		EmployeeNo:   user.EmployeeNo,
		IsRegistered: 1,
		Status:       *user.Status,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *H5UserService) UpdateUserInfo(userId uint, req request.UpdateUserInfoReq) (*response.H5UserResp, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	updates := make(map[string]interface{})

	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Phone != "" {
		// 检查手机号是否已被使用
		var existUser annual.AnnualUser
		result := global.GVA_DB.Where("phone = ? AND id != ?", req.Phone, userId).First(&existUser)
		if result.RowsAffected > 0 {
			return nil, errors.New("该手机号已被使用")
		}
		updates["phone"] = req.Phone
	}
	if req.Department != "" {
		updates["department"] = req.Department
	}
	if req.EmployeeNo != "" {
		// 检查工号是否已被使用
		var existUser annual.AnnualUser
		result := global.GVA_DB.Where("employee_no = ? AND id != ?", req.EmployeeNo, userId).First(&existUser)
		if result.RowsAffected > 0 {
			return nil, errors.New("该工号已被使用")
		}
		updates["employee_no"] = req.EmployeeNo
	}

	if len(updates) > 0 {
		if err := global.GVA_DB.Model(&user).Updates(updates).Error; err != nil {
			return nil, errors.New("更新失败")
		}
	}

	// 重新查询
	global.GVA_DB.First(&user, userId)

	isRegistered := 0
	if user.RealName != "" {
		isRegistered = 1
	}

	return &response.H5UserResp{
		ID:           user.ID,
		OpenId:       user.OpenId,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		RealName:     user.RealName,
		Phone:        user.Phone,
		Department:   user.Department,
		EmployeeNo:   user.EmployeeNo,
		IsRegistered: isRegistered,
		Status:       *user.Status,
	}, nil
}

// GetAuditStatus 获取审核状态
func (s *H5UserService) GetAuditStatus(userId uint) (*response.AuditStatusResp, error) {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &response.AuditStatusResp{
		Status:       *user.Status,
		RejectReason: user.RejectReason,
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
		ID:         user.ID,
		Nickname:   user.Nickname,
		Avatar:     user.Avatar,
		RealName:   user.RealName,
		Department: user.Department,
	}
}

// CheckUserRegistered 检查用户是否已报名并通过审核
func (s *H5UserService) CheckUserRegistered(userId uint) error {
	var user annual.AnnualUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.RealName == "" {
		return errors.New("请先完成报名")
	}

	if *user.Status == 0 {
		return errors.New("您的报名正在审核中")
	}

	if *user.Status == 2 {
		return errors.New("您的报名未通过审核")
	}

	return nil
}
