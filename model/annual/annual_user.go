package annual

import "lotteryBackend/global"

// AnnualUser 用户授权登录表
type AnnualUser struct {
	global.GVA_MODEL
	OpenId   string `json:"openid" gorm:"column:openid;type:varchar(64);uniqueIndex;not null;default:'';comment:微信openid"`
	UnionId  string `json:"unionid" gorm:"column:unionid;type:varchar(64);index;not null;default:'';comment:微信unionid"`
	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(64);not null;default:'';comment:微信昵称"`
	Avatar   string `json:"avatar" gorm:"column:avatar;type:varchar(512);not null;default:'';comment:头像URL"`
}

func (AnnualUser) TableName() string {
	return "annual_users"
}
