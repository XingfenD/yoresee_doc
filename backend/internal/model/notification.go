package model

import (
	"time"
)

type Notification struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID string    `gorm:"size:100;unique;index" json:"external_id"`
	ReceiverID int64     `gorm:"index;not null" json:"receiver_id"`
	Type       string    `gorm:"size:64;not null" json:"type"`
	Status     string    `gorm:"size:32;not null;default:unread" json:"status"`
	Title      string    `gorm:"size:255" json:"title"`
	Content    string    `gorm:"size:1024" json:"content"`
	Payload    string    `gorm:"type:jsonb" json:"payload"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
