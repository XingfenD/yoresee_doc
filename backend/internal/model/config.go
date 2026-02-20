package model

type SystemConfig struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Key   string `gorm:"size:100;unique;not null" json:"key"`
	Value string `gorm:"size:255;not null" json:"value"` // 'true'/'false' for boolean
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
