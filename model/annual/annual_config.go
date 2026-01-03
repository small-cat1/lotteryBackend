package annual

import "time"

// AnnualConfig 年会系统配置表
type AnnualConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey;comment:主键"`
	ConfigKey   string    `json:"configKey" gorm:"column:config_key;type:varchar(64);uniqueIndex;not null;default:'';comment:配置键"`
	ConfigValue string    `json:"configValue" gorm:"column:config_value;type:text;comment:配置值"`
	Description string    `json:"description" gorm:"column:description;type:varchar(255);not null;default:'';comment:配置说明"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;comment:创建时间"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;comment:更新时间"`
}

func (AnnualConfig) TableName() string {
	return "annual_configs"
}
