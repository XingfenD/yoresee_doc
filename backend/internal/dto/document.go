package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

type DocumentType string

const DocumentType_Markdown DocumentType = "markdown"

type DocumentBase struct {
	ExternalID string       `json:"external_id"`
	Title      string       `json:"title"`
	Type       DocumentType `json:"type"`
	Summary    string       `json:"summary"`
	Status     int          `json:"status"`
	Tags       []string     `json:"tags"`
	ViewCount  int          `json:"view_count"`
	EditCount  int          `json:"edit_count"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

type DocumentMetaResponse struct {
	DocumentBase
	Children    []*DocumentMetaResponse `json:"children,omitempty"`
	HasChildren bool                    `json:"hasChildren,omitempty"`
}

type DocumentResponse struct {
	DocumentMetaResponse
	Content string `json:"content"`
}

func NewDocumentMetaResponseFromModel(doc *model.Document) *DocumentMetaResponse {
	response := &DocumentMetaResponse{
		DocumentBase: DocumentBase{
			ExternalID: doc.ExternalID,
			Title:      doc.Title,
			Type:       DocumentType(doc.Type),
			Summary:    doc.Summary,
			Status:     doc.Status,
			Tags:       doc.Tags,
			ViewCount:  doc.ViewCount,
			EditCount:  doc.EditCount,
			CreatedAt:  doc.CreatedAt,
			UpdatedAt:  doc.UpdatedAt,
		},
	}

	return response
}

func NewDocumentResponseFromModel(doc *model.Document) *DocumentResponse {
	if doc == nil {
		return &DocumentResponse{}
	}
	response := &DocumentResponse{
		DocumentMetaResponse: DocumentMetaResponse{
			DocumentBase: DocumentBase{
				ExternalID: doc.ExternalID,
				Title:      doc.Title,
				Type:       DocumentType(doc.Type),
				Summary:    doc.Summary,
				Status:     doc.Status,
				Tags:       doc.Tags,
				ViewCount:  doc.ViewCount,
				EditCount:  doc.EditCount,
				CreatedAt:  doc.CreatedAt,
				UpdatedAt:  doc.UpdatedAt,
			},
		},
		Content: doc.Content,
	}

	return response
}

// type DirectoryResponse struct {
// 	ExternalID  string `json:"external_id"`
// 	Title       string `json:"title"`
// 	HasChildren bool   `json:"has_children"`
// 	ParentID    string `json:"parent_id"`
// }

type DocumentsListExternalArgs struct {
	UserExternalID         *string `json:"user_external_id"`
	RootDocumentExternalID *string `json:"root_document_external_id"`
	KnowledgeExternalID    *string `json:"knowledge_external_id"`
}

type DocumentsListFilterArgs struct {
	TitleKeyword         *string  `json:"title_keyword"`
	DocType              *string  `json:"doc_type"`
	Status               *int     `json:"status"`
	Tags                 []string `json:"tags"`
	CreateTimeRangeStart *string  `json:"create_time_range_start"`
	CreateTimeRangeEnd   *string  `json:"create_time_range_end"`
	UpdateTimeRangeStart *string  `json:"update_time_range_start"`
	UpdateTimeRangeEnd   *string  `json:"update_time_range_end"`
}

type ListDocumentsBaseArgs struct {
	ListOwnDoc    bool `json:"list_own_doc"`
	DirectoryOnly bool `json:"directory_only"`
}

type ListRecentDocumentsRequest struct {
	UserExternalID string     `json:"user_external_id"`
	StartTime      *time.Time `json:"start_time"`
	EndTime        *time.Time `json:"end_time"`
	Pagination     Pagination `json:"pagination"`
}

type RecordRecentDocumentRequest struct {
	UserExternalID     string `json:"user_external_id"`
	DocumentExternalID string `json:"document_external_id"`
}

type ListDocumentsByExternalReq struct {
	ExternalArgs *DocumentsListExternalArgs `json:"external_args"`
	ListDocumentsBaseArgs
	FilterArgs *DocumentsListFilterArgs `json:"filter_args"`
	SortArgs   SortArgs                 `json:"sort_args"`
	Pagination Pagination               `json:"pagination"`
	Options    *RecursiveOptions        `json:"options"`
}

type CreateDocumentReq struct {
	Title             string       `json:"title"`
	Type              DocumentType `json:"type"`
	CreatorExternalID *string      `json:"creator_external_id"`
	TemplateID        *int64       `json:"template_id"`

	// own doc or in knowledge_base
	CreateAsOwnDoc bool `json:"create_as_own_doc"`

	// optional
	ParentExternalID    *string `json:"parent_external_id"`
	KnowledgeExternalID *string `json:"knowledge_external_id"`
}

type CreateDocumentResponse struct {
	ExternalID string `json:"external_id"`
}

type UpdateDocumentRequest struct {
	ExternalID              string  `json:"external_id"`
	Title                   *string `json:"title,omitempty"`
	KnowledgeBaseExternalID *string `json:"knowledge_base_external_id,omitempty"`
	ParentExternalID        *string `json:"parent_external_id,omitempty"`
	MoveAsOwn               bool    `json:"move_as_own"`
	Content                 *string `json:"content,omitempty"`
}

type UpdateDocumentMetaRequest struct {
	ExternalID string    `json:"external_id"`
	Title      *string   `json:"title,omitempty"`
	Summary    *string   `json:"summary,omitempty"`
	Tags       *[]string `json:"tags,omitempty"`
	Status     *int      `json:"status,omitempty"`
}
