package model

import "time"

type DocumentVersion struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID    int64     `gorm:"not null;index" json:"document_id"`
	Version       int       `gorm:"not null" json:"version"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	ContentID     int64     `gorm:"not null;index" json:"content_id"`
	UserID        int64     `gorm:"not null;index" json:"user_id"`
	ChangeSummary string    `gorm:"type:text" json:"change_summary"`
	CreatedAt     time.Time `json:"created_at"`
}

func (DocumentVersion) TableName() string {
	return "document_versions"
}
