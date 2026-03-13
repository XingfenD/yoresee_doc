package service

import "github.com/XingfenD/yoresee_doc/internal/dto"

type documentsListMetaArgs struct {
	UserID      *int64 `json:"user_id"`
	ParentID    *int64 `json:"parent_id"`
	KnowledgeID *int64 `json:"knowledge_id"`
	ListOwnDoc  bool   `json:"list_own_doc"`
}

type documentsListReq struct {
	MetaArgs   *documentsListMetaArgs       `json:"meta_args"`
	FilterArgs *dto.DocumentsListFilterArgs `json:"filter_args"`
	SortArgs   dto.SortArgs                 `json:"sort_args"`
	Pagination dto.Pagination               `json:"pagination"`
	Options    *dto.RecursiveOptions        `json:"options"`
}

type knowledgeBaseListReq struct {
	CreatorID  *int64                           `json:"creator_id"`
	FilterArgs *dto.KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs   dto.SortArgs                     `json:"sort_args"`
	Pagination dto.Pagination                   `json:"pagination"`
}

type knowledgeBaseListOperation struct {
	req  *knowledgeBaseListReq
	srvc *KnowledgeBaseService

	withDocumentExtend bool
	withUserExtend     bool
}
