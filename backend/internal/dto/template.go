package dto

import "time"

type TemplateContainer string

const (
	TemplateContainerOwn          TemplateContainer = "own"
	TemplateContainerKnowledgeBase TemplateContainer = "knowledge_base"
	TemplateContainerPublic        TemplateContainer = "public"
)

type CreateTemplateRequest struct {
	UserExternalID          string            `json:"user_external_id"`
	TargetContainer         TemplateContainer `json:"target_container"`
	KnowledgeBaseExternalID *string           `json:"knowledge_base_external_id,omitempty"`
	TemplateContent         string            `json:"template_content"`
	Type                    DocumentType      `json:"type"`
}

type UpdateTemplateSettingsRequest struct {
	UserExternalID string  `json:"user_external_id"`
	TemplateID     int64   `json:"template_id"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	IsPublic       *bool   `json:"is_public,omitempty"`
}

type TemplateResponse struct {
	ID                      int64        `json:"id"`
	Name                    string       `json:"name"`
	Description             string       `json:"description"`
	Content                 string       `json:"content"`
	Type                    DocumentType `json:"type"`
	Scope                   string       `json:"scope"`
	KnowledgeBaseExternalID string       `json:"knowledge_base_external_id"`
	Tags                    []string     `json:"tags"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at"`
}

type TemplateListFilterArgs struct {
	NameKeyword     *string            `json:"name_keyword"`
	TargetContainer *TemplateContainer `json:"target_container"`
	KnowledgeBaseID *string            `json:"knowledge_base_id"`
	Type            *DocumentType      `json:"type"`
}

type TemplateListByExternalRequest struct {
	CreatorExternalID string                  `json:"creator_external_id"`
	FilterArgs        *TemplateListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                `json:"sort_args"`
	Pagination        Pagination              `json:"pagination"`
}

type CreateRecentTemplateRequest struct {
	UserExternalID string    `json:"user_external_id"`
	TemplateID     int64     `json:"template_id"`
	AccessTime     time.Time `json:"access_time"`
}

type ListRecentTemplatesRequest struct {
	UserExternalID string     `json:"user_external_id"`
	StartTime      *time.Time `json:"start_time,omitempty"`
	EndTime        *time.Time `json:"end_time,omitempty"`
	Pagination     Pagination `json:"pagination"`
}
