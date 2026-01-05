package annual

import "time"

// AnnualCheckIn 年会签到表（包含个人信息）
type AnnualCheckIn struct {
	ID           uint      `json:"id" gorm:"primaryKey;comment:主键"`
	ActivityId   uint      `json:"activityId" gorm:"column:activity_id;not null;default:0;uniqueIndex:uk_activity_user;comment:活动ID"`
	UserId       uint      `json:"userId" gorm:"column:user_id;index;not null;default:0;uniqueIndex:uk_activity_user;comment:用户ID"`
	RealName     string    `json:"realName" gorm:"column:real_name;type:varchar(32);not null;default:'';comment:真实姓名"`
	Phone        string    `json:"phone" gorm:"column:phone;type:varchar(20);index;not null;default:'';comment:手机号"`
	Department   string    `json:"department" gorm:"column:department;type:varchar(64);not null;default:'';comment:部门"`
	EmployeeNo   string    `json:"employeeNo" gorm:"column:employee_no;type:varchar(32);not null;default:'';comment:工号"`
	Status       int       `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:0;comment:状态：0待审核 1已通过 2已拒绝"`
	RejectReason string    `json:"rejectReason" gorm:"column:reject_reason;type:varchar(255);not null;default:'';comment:拒绝原因"`
	CheckInTime  time.Time `json:"checkInTime" gorm:"column:check_in_time;not null;default:CURRENT_TIMESTAMP;comment:签到时间"`
	Ip           string    `json:"ip" gorm:"column:ip;type:varchar(64);not null;default:'';comment:签到IP"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	// 关联
	User     AnnualUser     `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Activity AnnualActivity `json:"activity" gorm:"foreignKey:ActivityId;references:ID"`
}

func (AnnualCheckIn) TableName() string {
	return "annual_check_ins"
}
