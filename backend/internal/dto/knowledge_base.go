package dto

import "time"

type KnowledgeBaseBase struct {
	ExternalID    string    `json:"external_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Cover         string    `json:"cover"`
	CreatorUserID int64     `json:"creator_user_id"`
	IsPublic      bool      `json:"is_public"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
