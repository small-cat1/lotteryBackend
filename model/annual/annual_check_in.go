package annual

import "time"

// AnnualCheckIn 年会签到表
type AnnualCheckIn struct {
	ID          uint      `json:"id" gorm:"primaryKey;comment:主键"`
	ActivityId  uint      `json:"activityId" gorm:"column:activity_id;not null;default:0;uniqueIndex:uk_activity_user;comment:活动ID"`
	UserId      uint      `json:"userId" gorm:"column:user_id;index;not null;default:0;uniqueIndex:uk_activity_user;comment:用户ID"`
	CheckInTime time.Time `json:"checkInTime" gorm:"column:check_in_time;not null;default:CURRENT_TIMESTAMP;comment:签到时间"`
	Ip          string    `json:"ip" gorm:"column:ip;type:varchar(64);not null;default:'';comment:签到IP"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	// 关联
	User     AnnualUser     `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Activity AnnualActivity `json:"activity" gorm:"foreignKey:ActivityId;references:ID"`
}

func (AnnualCheckIn) TableName() string {
	return "annual_check_ins"
}
