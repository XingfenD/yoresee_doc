package model

import (
	"time"

	"gorm.io/gorm"
)

type DocumentType string

const DocumentType_Markdown DocumentType = "markdown"

type ContainerType string

const ContainerType_Own ContainerType = "own"
const ContainerType_KnowledgeBase ContainerType = "knowledge_base"

type Document struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID    string         `gorm:"size:100;unique;index;not null" json:"external_id"`
	Title         string         `gorm:"size:200;not null" json:"title"`
	Type          DocumentType   `gorm:"size:20;default:'markdown'" json:"type"`
	Summary       string         `gorm:"type:text" json:"summary"`
	ParentID      int64          `gorm:"default:0;index" json:"parent_id"` // 0 means root
	UserID        int64          `gorm:"not null;index" json:"user_id"`
	KnowledgeID   *int64         `gorm:"index" json:"knowledge_id"`
	ContainerType ContainerType  `gorm:"not null;index" json:"containter_type"`
	IsPublic      bool           `gorm:"default:false;index" json:"is_public"`
	Tags          []string       `gorm:"serializer:json" json:"tags"`
	Path          string         `gorm:"type:ltree;not null;index:idx_path_gist,using:gist"`
	Depth         int            `gorm:"not null;index"`
	ViewCount     int            `gorm:"default:0" json:"view_count"`
	EditCount     int            `gorm:"default:0" json:"edit_count"`
	Content       string         `gorm:"type:text" json:"content"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Document) TableName() string {
	return "document_metas"
}
