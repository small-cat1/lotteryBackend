package annual

import (
	"lotteryBackend/global"
	"time"
)

// AnnualShakeRound 摇一摇游戏场次表
type AnnualShakeRound struct {
	global.GVA_MODEL
	ActivityId  uint       `json:"activityId" gorm:"column:activity_id;index;not null;default:0;comment:活动ID"`
	RoundName   string     `json:"roundName" gorm:"column:round_name;type:varchar(64);not null;default:'';comment:场次名称"`
	Duration    int        `json:"duration" gorm:"column:duration;type:int unsigned;not null;default:30;comment:游戏时长(秒)"`
	WinnerCount int        `json:"winnerCount" gorm:"column:winner_count;type:int unsigned;not null;default:1;comment:本轮获奖人数"`
	PrizeId     uint       `json:"prizeId" gorm:"column:prize_id;not null;default:0;comment:关联奖品ID"`
	Status      *int       `json:"status" gorm:"column:status;type:tinyint unsigned;index;not null;default:0;comment:状态：0未开始 1进行中 2已结束"`
	StartTime   *time.Time `json:"startTime" gorm:"column:start_time;comment:实际开始时间"`
	EndTime     *time.Time `json:"endTime" gorm:"column:end_time;comment:实际结束时间"`
	Sort        int        `json:"sort" gorm:"column:sort;type:int unsigned;not null;default:0;comment:排序"`
	// 关联
	Activity AnnualActivity `json:"activity" gorm:"foreignKey:ActivityId;references:ID"`
	Prize    AnnualPrize    `json:"prize" gorm:"foreignKey:PrizeId;references:ID"`
}

func (AnnualShakeRound) TableName() string {
	return "annual_shake_rounds"
}
