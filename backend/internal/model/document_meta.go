package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentMeta struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID  string         `gorm:"size:100;unique;not null" json:"external_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Type        string         `gorm:"size:20;default:'markdown'" json:"type"`
	Summary     string         `gorm:"type:text" json:"summary"`
	ParentID    int64          `gorm:"default:0;index" json:"parent_id"` // 0 means root
	UserID      int64          `gorm:"not null;index" json:"user_id"`
	KnowledgeID *int64         `gorm:"index" json:"knowledge_id"`
	Status      int            `gorm:"default:1" json:"status"`
	IsPublic    bool           `gorm:"default:false" json:"is_public"`
	Tags        []string       `gorm:"serializer:json" json:"tags"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	EditCount   int            `gorm:"default:0" json:"edit_count"`
	Version     int            `gorm:"default:1" json:"version"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// virtual field for children (used in API response)
	Children    []*DocumentMeta `gorm:"-" json:"children,omitempty"`
	HasChildren bool            `gorm:"-" json:"hasChildren,omitempty"` // For lazy loading
}

func (DocumentMeta) TableName() string {
	return "document_metas"
}
