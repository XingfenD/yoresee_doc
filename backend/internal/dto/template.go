package dto

import "time"

type TemplateContainer int

const (
	TemplateContainerOwn TemplateContainer = iota
	TemplateContainerKnowledgeBase
	TemplateContainerPublic
)

type CreateTemplateRequest struct {
	UserExternalID          string
	TargetContainer         TemplateContainer
	KnowledgeBaseExternalID *string
	TemplateContent         string
}

type TemplateResponse struct {
	ID                      int64     `json:"id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	Content                 string    `json:"content"`
	Scope                   string    `json:"scope"`
	KnowledgeBaseExternalID string    `json:"knowledge_base_external_id"`
	Tags                    []string  `json:"tags"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type TemplateListFilterArgs struct {
	NameKeyword     *string            `json:"name_keyword"`
	TargetContainer *TemplateContainer `json:"target_container"`
	KnowledgeBaseID *string            `json:"knowledge_base_id"`
}

type TemplateListByExternalReq struct {
	CreatorExternalID string                  `json:"creator_external_id"`
	FilterArgs        *TemplateListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                `json:"sort_args"`
	Pagination        Pagination              `json:"pagination"`
}

type CreateRecentTemplateRequest struct {
	UserExternalID string
	TemplateID     int64
	AccessTime     time.Time
}

type ListRecentTemplatesRequest struct {
	UserExternalID string
	StartTime      *time.Time
	EndTime        *time.Time
	Pagination     Pagination
}
