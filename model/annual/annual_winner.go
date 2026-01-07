package annual

import "time"

// AnnualWinner 年会中奖记录表
type AnnualWinner struct {
	ID          uint       `json:"id" gorm:"primaryKey;comment:主键"`
	ActivityId  uint       `json:"activityId" gorm:"column:activity_id;index;not null;default:0;comment:活动ID"`
	UserId      uint       `json:"userId" gorm:"column:user_id;index;not null;default:0;comment:用户ID"`
	PrizeId     uint       `json:"prizeId" gorm:"column:prize_id;index;not null;default:0;comment:奖品ID"`
	RoundId     uint       `json:"roundId" gorm:"column:round_id;index;not null;default:0;comment:摇一摇场次ID，普通抽奖为0"`
	WinType     *int       `json:"winType" gorm:"column:win_type;type:tinyint unsigned;not null;default:1;comment:中奖方式：1摇一摇 2随机抽奖 3弹幕抽奖"`
	Status      *int       `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:0;comment:领奖状态：0未领取 1已领取"`
	ReceiveCode string     `json:"receiveCode" gorm:"column:receive_code;type:varchar(10);not null;default:'';uniqueIndex;comment:核销密码"` // ✅ 新增
	ReceiveTime *time.Time `json:"receiveTime" gorm:"column:receive_time;comment:领奖时间"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;comment:中奖时间"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	// 关联
	User     AnnualUser       `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Prize    AnnualPrize      `json:"prize" gorm:"foreignKey:PrizeId;references:ID"`
	Activity AnnualActivity   `json:"activity" gorm:"foreignKey:ActivityId;references:ID"`
	Round    AnnualShakeRound `json:"round" gorm:"foreignKey:RoundId;references:ID"`

	// ✅ 新增：非数据库字段，用于返回成绩信息
	Score int `json:"score" gorm:"-"`
	Rank  int `json:"rank" gorm:"-"`
}

func (AnnualWinner) TableName() string {
	return "annual_winners"
}
