package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

type KnowledgeBaseBase struct {
	ExternalID  string    `json:"external_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cover       string    `json:"cover"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type KnowledgeBaseExtend struct {
	CreatorUserExternalID string `json:"creator_user_external_id"`
	CreatorName           string `json:"creator_name"`
	DocumentsCount        int64  `json:"documents_count"`
}

type KnowledgeBaseResponse struct {
	KnowledgeBaseBase
	KnowledgeBaseExtend
}

func NewKnowledgeBaseResponseFromModel(kb *model.KnowledgeBase, kbExtend *KnowledgeBaseExtend) *KnowledgeBaseResponse {
	response := &KnowledgeBaseResponse{
		KnowledgeBaseBase: KnowledgeBaseBase{
			ExternalID:  kb.ExternalID,
			Name:        kb.Name,
			Description: kb.Description,
			Cover:       kb.Cover,
			IsPublic:    kb.IsPublic,
			CreatedAt:   kb.CreatedAt,
			UpdatedAt:   kb.UpdatedAt,
			DeletedAt:   kb.DeletedAt,
		},
	}
	if kbExtend != nil {
		response.KnowledgeBaseExtend = KnowledgeBaseExtend{
			CreatorUserExternalID: kbExtend.CreatorUserExternalID,
			CreatorName:           kbExtend.CreatorName,
			DocumentsCount:        kbExtend.DocumentsCount,
		}
	}

	return response
}

type CreateRecentKnowledgeBaseRequest struct {
	UserExternalID          string
	KnowledgeBaseExternalID string
	AssessTime              time.Time
}
