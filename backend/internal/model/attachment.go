package model

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID   string         `gorm:"size:100;unique;not null" json:"external_id"`
	DocumentID   int64          `gorm:"not null;index" json:"document_id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Path         string         `gorm:"size:512;not null" json:"path"`
	Size         int64          `json:"size"`
	MimeType     string         `gorm:"size:100" json:"mime_type"`
	UserID       int64          `gorm:"not null" json:"user_id"`
	PresignedURL string         `gorm:"-" json:"presigned_url"` // 预签名URL，不存储到数据库
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Attachment) TableName() string {
	return "attachments"
}
