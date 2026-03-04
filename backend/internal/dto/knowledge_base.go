package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

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

type KnowledgeBaseResponse struct {
	KnowledgeBaseBase
	CreatorName    string `json:"creator_name"`
	DocumentsCount int    `json:"documents_count"`
}

func NewKnowledgeBaseResponseFromModel(kb *model.KnowledgeBase) *KnowledgeBaseResponse {
	response := &KnowledgeBaseResponse{
		KnowledgeBaseBase: KnowledgeBaseBase{
			ExternalID:    kb.ExternalID,
			Name:          kb.Name,
			Description:   kb.Description,
			Cover:         kb.Cover,
			CreatorUserID: kb.CreatorUserID,
			IsPublic:      kb.IsPublic,
			CreatedAt:     kb.CreatedAt,
			UpdatedAt:     kb.UpdatedAt,
			DeletedAt:     kb.DeletedAt,
		},
		CreatorName:    "", // Will be populated later
		DocumentsCount: 0,  // Will be populated later
	}

	return response
}

type CreateRecentKnowledgeBaseRequest struct {
	UserExternalID          string
	KnowledgeBaseExternalID string
	AssessTime              time.Time
}
