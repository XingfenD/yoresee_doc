package dto

import "github.com/XingfenD/yoresee_doc/internal/dto"

type DocumentsListMetaArgs struct {
	UserID      *int64 `json:"user_id"`
	ParentID    *int64 `json:"parent_id"`
	KnowledgeID *int64 `json:"knowledge_id"`
}

type DocumentsListReq struct {
	MetaArgs *DocumentsListMetaArgs `json:"meta_args"`
	dto.ListDocumentsBaseArgs
	FilterArgs *dto.DocumentsListFilterArgs `json:"filter_args"`
	SortArgs   dto.SortArgs                 `json:"sort_args"`
	Pagination dto.Pagination               `json:"pagination"`
	Options    *dto.RecursiveOptions        `json:"options"`
}

type KnowledgeBaseListReq struct {
	CreatorID  *int64                           `json:"creator_id"`
	FilterArgs *dto.KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs   dto.SortArgs                     `json:"sort_args"`
	Pagination dto.Pagination                   `json:"pagination"`
}
