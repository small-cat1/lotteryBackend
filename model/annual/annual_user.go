package annual

import "lotteryBackend/global"

// AnnualUser 年会用户表
type AnnualUser struct {
	global.GVA_MODEL
	Openid       string `json:"openid" gorm:"column:openid;type:varchar(64);uniqueIndex;not null;default:'';comment:微信openid"`
	Unionid      string `json:"unionid" gorm:"column:unionid;type:varchar(64);index;not null;default:'';comment:微信unionid"`
	Nickname     string `json:"nickname" gorm:"column:nickname;type:varchar(64);not null;default:'';comment:微信昵称"`
	Avatar       string `json:"avatar" gorm:"column:avatar;type:varchar(512);not null;default:'';comment:头像URL"`
	IsRegistered *int   `json:"isRegistered" gorm:"column:is_registered;type:tinyint unsigned;not null;default:0;comment:是否已报名：0游客 1已报名"`
	RealName     string `json:"realName" gorm:"column:real_name;type:varchar(32);not null;default:'';comment:真实姓名"`
	Department   string `json:"department" gorm:"column:department;type:varchar(64);not null;default:'';comment:部门"`
	Phone        string `json:"phone" gorm:"column:phone;type:varchar(20);index;not null;default:'';comment:手机号"`
	EmployeeNo   string `json:"employeeNo" gorm:"column:employee_no;type:varchar(32);not null;default:'';comment:工号"`
	Status       *int   `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:状态：0待审核 1已通过 2已拒绝"`
}

func (AnnualUser) TableName() string {
	return "annual_users"
}
