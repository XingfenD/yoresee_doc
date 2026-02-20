package model

import (
	"time"
)

type TemplateMeta struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	Description     string    `gorm:"type:text" json:"description"`
	ContentID       int64     `gorm:"not null;index" json:"content_id"`
	UserID          int64     `gorm:"not null" json:"user_id"`
	CategoryID      int64     `json:"category_id"`
	Scope           string    `gorm:"size:20;default:'private'" json:"scope"` // private, system, knowledge_base, shared
	KnowledgeBaseID *int64    `json:"knowledge_base_id"`
	IsPublic        bool      `gorm:"default:false" json:"is_public"` // Deprecated, keep for migration
	Tags            []string  `gorm:"serializer:json" json:"tags"`
	UsageCount      int       `gorm:"default:0" json:"usage_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (TemplateMeta) TableName() string {
	return "templates"
}
