package annual

import (
	"errors"
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
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id DESC").Find(&list).Error
	return
}

// DeleteUser 删除用户（安全检查）
func (s *AnnualUserService) DeleteUser(id uint) (err error) {
	if err = s.checkUserRelations(id); err != nil {
		return err
	}
	return global.GVA_DB.Delete(&annual.AnnualUser{}, id).Error
}

// DeleteUserByIds 批量删除用户（安全检查）
func (s *AnnualUserService) DeleteUserByIds(ids []uint) (err error) {
	for _, id := range ids {
		if err = s.checkUserRelations(id); err != nil {
			return err
		}
	}
	return global.GVA_DB.Delete(&annual.AnnualUser{}, ids).Error
}

// checkUserRelations 检查用户关联数据
func (s *AnnualUserService) checkUserRelations(userId uint) error {
	var count int64

	// 检查签到记录
	if err := global.GVA_DB.Model(&annual.AnnualCheckIn{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该用户存在签到记录，无法删除")
	}

	// 检查弹幕记录
	if err := global.GVA_DB.Model(&annual.AnnualDanmaku{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该用户存在弹幕记录，无法删除")
	}

	// 检查摇一摇成绩
	if err := global.GVA_DB.Model(&annual.AnnualShakeScore{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该用户存在摇一摇成绩，无法删除")
	}

	// 检查中奖记录
	if err := global.GVA_DB.Model(&annual.AnnualWinner{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该用户存在中奖记录，无法删除")
	}

	return nil
}
