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

type CreateKnowledgeBaseRequest struct {
	CreatorExternalID string `json:"creator_external_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Cover             string `json:"cover"`
	IsPublic          bool   `json:"is_public"`
}

type CreateKnowledgeBaseResponse struct {
	ExternalID string `json:"external_id"`
}

type ListRecentKnowledgeBasesRequest struct {
	UserExternalID string
	StartTime      *time.Time
	EndTime        *time.Time
	Pagination     Pagination
}

type KnowledgeBaseListByExternalReq struct {
	CreatorExternalID string                       `json:"creator_external_id"`
	FilterArgs        *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                     `json:"sort_args"`
	Pagination        Pagination                   `json:"pagination"`
}

type KnowledgeBaseGetByExternalIDReq struct {
	KnowledgeBaseExternalID string `json:"knowledge_base_external_id"`
}

type KnowledgeBaseListFilterArgs struct {
	IsPublic             *bool   `json:"is_public"`
	NameKeyword          *string `json:"name_keyword"`
	CreateTimeRangeStart *string `json:"create_time_range_start"`
	CreateTimeRangeEnd   *string `json:"create_time_range_end"`
	UpdateTimeRangeStart *string `json:"update_time_range_start"`
	UpdateTimeRangeEnd   *string `json:"update_time_range_end"`
}
