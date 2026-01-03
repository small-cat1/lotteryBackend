package annual

import (
	"lotteryBackend/global"
	"time"
)

// AnnualActivity 年会活动表
type AnnualActivity struct {
	global.GVA_MODEL
	Title          string     `json:"title" gorm:"column:title;type:varchar(128);not null;default:'';comment:活动名称"`
	Cover          string     `json:"cover" gorm:"column:cover;type:varchar(512);not null;default:'';comment:活动封面图"`
	Description    string     `json:"description" gorm:"column:description;type:text;comment:活动描述"`
	StartTime      *time.Time `json:"startTime" gorm:"column:start_time;comment:活动开始时间"`
	EndTime        *time.Time `json:"endTime" gorm:"column:end_time;comment:活动结束时间"`
	CheckInEnabled *int       `json:"checkInEnabled" gorm:"column:check_in_enabled;type:tinyint unsigned;not null;default:1;comment:是否开启签到：0否 1是"`
	DanmakuEnabled *int       `json:"danmakuEnabled" gorm:"column:danmaku_enabled;type:tinyint unsigned;not null;default:1;comment:是否开启弹幕：0否 1是"`
	DanmakuAudit   *int       `json:"danmakuAudit" gorm:"column:danmaku_audit;type:tinyint unsigned;not null;default:0;comment:弹幕是否需要审核：0否 1是"`
	WinnerExclude  *int       `json:"winnerExclude" gorm:"column:winner_exclude;type:tinyint unsigned;not null;default:0;comment:中奖后排除后续抽奖：0否 1是"`
	Status         *int       `json:"status" gorm:"column:status;type:tinyint unsigned;index;not null;default:0;comment:状态：0未开始 1进行中 2已结束"`
	CreatedBy      uint       `json:"createdBy" gorm:"column:created_by;not null;default:0;comment:创建人"`
}

func (AnnualActivity) TableName() string {
	return "annual_activities"
}
