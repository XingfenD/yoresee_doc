package model

import "time"

type RecentDocument struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	DocumentID int64     `gorm:"not null;index" json:"document_id"`
	AccessedAt time.Time `gorm:"not null" json:"accessed_at"`
}

func (RecentDocument) TableName() string {
	return "recent_documents"
}
