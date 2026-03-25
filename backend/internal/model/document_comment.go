package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentComment struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID string         `gorm:"size:100;unique;index;not null" json:"external_id"`
	DocumentID int64          `gorm:"index;not null" json:"document_id"`
	ParentID   int64          `gorm:"default:0;index" json:"parent_id"`
	CreatorID  int64          `gorm:"index;not null" json:"creator_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (DocumentComment) TableName() string {
	return "document_comments"
}
