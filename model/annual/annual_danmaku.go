package annual

import "time"

// AnnualDanmaku 年会弹幕表
type AnnualDanmaku struct {
	ID         uint      `json:"id" gorm:"primaryKey;comment:主键"`
	ActivityId uint      `json:"activityId" gorm:"column:activity_id;index:idx_activity_status;not null;default:0;comment:活动ID"`
	UserId     uint      `json:"userId" gorm:"column:user_id;index;not null;default:0;comment:用户ID"`
	Content    string    `json:"content" gorm:"column:content;type:varchar(200);not null;default:'';comment:弹幕内容"`
	Color      string    `json:"color" gorm:"column:color;type:varchar(16);not null;default:'#FFFFFF';comment:弹幕颜色"`
	Status     *int      `json:"status" gorm:"column:status;type:tinyint unsigned;index:idx_activity_status;not null;default:1;comment:状态：0待审核 1已通过 2已拒绝"`
	IsTop      *int      `json:"isTop" gorm:"column:is_top;type:tinyint unsigned;not null;default:0;comment:是否置顶：0否 1是"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;index;comment:创建时间"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	// 关联
	User AnnualUser `json:"user" gorm:"foreignKey:UserId;references:ID"`
}

func (AnnualDanmaku) TableName() string {
	return "annual_danmaku"
}
