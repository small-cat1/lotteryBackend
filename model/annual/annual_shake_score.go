package annual

import "time"

// AnnualShakeScore 摇一摇成绩表
type AnnualShakeScore struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
	RoundId   uint      `json:"roundId" gorm:"column:round_id;not null;default:0;uniqueIndex:uk_round_user;index:idx_score;comment:场次ID"`
	UserId    uint      `json:"userId" gorm:"column:user_id;index;not null;default:0;uniqueIndex:uk_round_user;comment:用户ID"`
	Score     int       `json:"score" gorm:"column:score;type:int unsigned;not null;default:0;comment:摇动次数/得分"`
	Rank      int       `json:"rank" gorm:"column:rank;type:int unsigned;not null;default:0;comment:最终排名"`
	IsWinner  *int      `json:"isWinner" gorm:"column:is_winner;type:tinyint unsigned;not null;default:0;comment:是否获奖：0否 1是"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
	// 关联
	User  AnnualUser       `json:"user" gorm:"foreignKey:UserId;references:ID"`
	Round AnnualShakeRound `json:"round" gorm:"foreignKey:RoundId;references:ID"`
}

func (AnnualShakeScore) TableName() string {
	return "annual_shake_scores"
}
