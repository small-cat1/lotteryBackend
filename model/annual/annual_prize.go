package annual

import "lotteryBackend/global"

// AnnualPrize 年会奖品表
type AnnualPrize struct {
	global.GVA_MODEL
	ActivityId  uint   `json:"activityId" gorm:"column:activity_id;index;not null;default:0;comment:活动ID"`
	Name        string `json:"name" gorm:"column:name;type:varchar(64);not null;default:'';comment:奖品名称"`
	Image       string `json:"image" gorm:"column:image;type:varchar(512);not null;default:'';comment:奖品图片"`
	Level       *int   `json:"level" gorm:"column:level;type:tinyint unsigned;index;not null;default:1;comment:奖品等级：1特等奖 2一等奖 3二等奖 4三等奖 5参与奖"`
	TotalCount  int    `json:"totalCount" gorm:"column:total_count;type:int unsigned;not null;default:0;comment:奖品总数"`
	RemainCount int    `json:"remainCount" gorm:"column:remain_count;type:int unsigned;not null;default:0;comment:剩余数量"`
	Sort        int    `json:"sort" gorm:"column:sort;type:int unsigned;not null;default:0;comment:排序"`
}

func (AnnualPrize) TableName() string {
	return "annual_prizes"
}
