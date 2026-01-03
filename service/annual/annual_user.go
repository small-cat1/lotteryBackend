package annual

import (
	"lotteryBackend/global"
	"lotteryBackend/model/annual"
	annualReq "lotteryBackend/model/annual/request"
)

type AnnualUserService struct{}

// GetUserList 获取用户列表
func (s *AnnualUserService) GetUserList(info annualReq.UserSearch) (list []annual.AnnualUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&annual.AnnualUser{})

	if info.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+info.Nickname+"%")
	}
	if info.RealName != "" {
		db = db.Where("real_name LIKE ?", "%"+info.RealName+"%")
	}
	if info.Department != "" {
		db = db.Where("department LIKE ?", "%"+info.Department+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.IsRegistered != nil {
		db = db.Where("is_registered = ?", *info.IsRegistered)
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

// GetUserById 获取用户详情
func (s *AnnualUserService) GetUserById(id string) (user annual.AnnualUser, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&user).Error
	return
}

// UpdateUser 更新用户
func (s *AnnualUserService) UpdateUser(user annual.AnnualUser) (err error) {
	return global.GVA_DB.Model(&annual.AnnualUser{}).Where("id = ?", user.ID).Updates(&user).Error
}

// UpdateUserStatus 更新用户状态
func (s *AnnualUserService) UpdateUserStatus(id uint, status int) (err error) {
	return global.GVA_DB.Model(&annual.AnnualUser{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteUser 删除用户
func (s *AnnualUserService) DeleteUser(id uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualUser{}, id).Error
}

// DeleteUserByIds 批量删除用户
func (s *AnnualUserService) DeleteUserByIds(ids []uint) (err error) {
	return global.GVA_DB.Delete(&annual.AnnualUser{}, ids).Error
}

// ExportUser 导出用户
func (s *AnnualUserService) ExportUser(info annualReq.UserSearch) (filePath string, err error) {
	// TODO: 实现导出逻辑
	return "", nil
}
