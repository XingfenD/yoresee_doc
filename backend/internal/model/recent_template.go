package model

import "time"

type RecentTemplate struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	TemplateID int64     `gorm:"not null;index" json:"template_id"`
	AccessedAt time.Time `gorm:"not null" json:"accessed_at"`
}

func (RecentTemplate) TableName() string {
	return "recent_templates"
}
