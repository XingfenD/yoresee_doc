package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

type DocumentBase struct {
	ExternalID string    `json:"external_id"`
	Title      string    `json:"title"`
	Type       string    `json:"type"`
	Summary    string    `json:"summary"`
	Status     int       `json:"status"`
	Tags       []string  `json:"tags"`
	ViewCount  int       `json:"view_count"`
	EditCount  int       `json:"edit_count"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentResponse struct {
	DocumentBase
	Children    []*DocumentResponse `json:"children,omitempty"`
	HasChildren bool                `json:"hasChildren,omitempty"`
}

func NewDocumentResponseFromModel(doc *model.DocumentMeta) *DocumentResponse {
	response := &DocumentResponse{
		DocumentBase: DocumentBase{
			ExternalID: doc.ExternalID,
			Title:      doc.Title,
			Type:       doc.Type,
			Summary:    doc.Summary,
			Status:     doc.Status,
			Tags:       doc.Tags,
			ViewCount:  doc.ViewCount,
			EditCount:  doc.EditCount,
			Version:    doc.Version,
			CreatedAt:  doc.CreatedAt,
			UpdatedAt:  doc.UpdatedAt,
		},
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

type ListDocumentsByExternalReq struct {
	ExternalArgs *DocumentsListExternalArgs `json:"external_args"`
	FilterArgs   *DocumentsListFilterArgs   `json:"filter_args"`
	SortArgs     SortArgs                   `json:"sort_args"`
	Pagination   Pagination                 `json:"pagination"`
	Options      *RecursiveOptions          `json:"options"`
}

type CreateDocumentReq struct {
	Title             string   `json:"title"`
	Type              string   `json:"type"`
	Summary           string   `json:"summary"`
	Status            int      `json:"status"`
	Tags              []string `json:"tags"`
	Content           string   `json:"content"`
	CreatorExternalID *string  `json:"creator_external_id"`

	//
	CreateAsOwnDoc bool `json:"create_as_own_doc"`

	// optional
	ParentExternalID    *string `json:"parent_external_id"`
	KnowledgeExternalID *string `json:"knowledge_external_id"`
}

func (req *CreateDocumentReq) Validate() error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.CreateAsOwnDoc && req.KnowledgeExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID not nil when CreateAsOwnDoc")
	}
	if !req.CreateAsOwnDoc && req.KnowledgeExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID is nil when not CreateAsOwnDoc")
	}

	return nil
}

type CreateDocumentResponse struct {
	ExternalID string `json:"external_id"`
}
