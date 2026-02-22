package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

type DocumentBase struct {
	ID          int64     `json:"id"`
	ExternalID  string    `json:"external_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Summary     string    `json:"summary"`
	ParentID    int64     `json:"parent_id"`
	UserID      int64     `json:"user_id"`
	KnowledgeID *int64    `json:"knowledge_id"`
	Status      int       `json:"status"`
	IsPublic    bool      `json:"is_public"`
	Tags        []string  `json:"tags"`
	ViewCount   int       `json:"view_count"`
	EditCount   int       `json:"edit_count"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DocumentResponse struct {
	DocumentBase
	Children    []*DocumentResponse `json:"children,omitempty"`
	HasChildren bool                `json:"hasChildren,omitempty"`
}

func NewDocumentResponseFromModel(doc *model.DocumentMeta) *DocumentResponse {
	response := &DocumentResponse{
		DocumentBase: DocumentBase{
			ID:          doc.ID,
			ExternalID:  doc.ExternalID,
			Title:       doc.Title,
			Type:        doc.Type,
			Summary:     doc.Summary,
			ParentID:    doc.ParentID,
			UserID:      doc.UserID,
			KnowledgeID: doc.KnowledgeID,
			Status:      doc.Status,
			IsPublic:    doc.IsPublic,
			Tags:        doc.Tags,
			ViewCount:   doc.ViewCount,
			EditCount:   doc.EditCount,
			Version:     doc.Version,
			CreatedAt:   doc.CreatedAt,
			UpdatedAt:   doc.UpdatedAt,
		},
		HasChildren: doc.HasChildren,
	}

	// 递归转换子文档
	if len(doc.Children) > 0 {
		response.Children = make([]*DocumentResponse, len(doc.Children))
		for i, child := range doc.Children {
			response.Children[i] = NewDocumentResponseFromModel(child)
		}
	}

	return response
}
