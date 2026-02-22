package model

import (
	"time"
)

type KnowledgeBase struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID     string    `gorm:"size:100;unique;not null;index" json:"external_id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	Description    string    `gorm:"size:255" json:"description"`
	Cover          string    `gorm:"size:255" json:"cover"`
	UserID         int64     `gorm:"not null" json:"user_id"` // owner user id
	IsPublic       bool      `gorm:"default:false" json:"is_public"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastAccessedAt time.Time `gorm:"-" json:"last_accessed_at"` // Auxiliary for recent
}

type RecentKnowledgeBase struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"not null;index" json:"user_id"`
	KnowledgeBaseID int64     `gorm:"not null;index" json:"knowledge_base_id"`
	AccessedAt      time.Time `gorm:"not null" json:"accessed_at"`
}

func (RecentKnowledgeBase) TableName() string {
	return "recent_knowledge_bases"
}

func (KnowledgeBase) TableName() string {
	return "knowledge_bases"
}
