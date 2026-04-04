package model

import (
	"time"
)

type Template struct {
	ID              int64        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string       `gorm:"size:100;not null" json:"name"`
	Description     string       `gorm:"type:text" json:"description"`
	DocumentType    DocumentType `gorm:"column:document_type;size:20;default:'markdown'" json:"document_type"`
	Content         string       `gorm:"type:text" json:"content"`
	UserID          int64        `gorm:"not null" json:"user_id"`
	Scope           string       `gorm:"size:20;default:'private'" json:"scope"` // private, system, knowledge_base
	KnowledgeBaseID *int64       `json:"knowledge_base_id"`
	IsPublic        bool         `gorm:"default:false" json:"is_public"` // Deprecated, keep for migration
	Tags            []string     `gorm:"serializer:json" json:"tags"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func (Template) TableName() string {
	return "templates"
}
