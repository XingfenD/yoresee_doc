package service

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
)

type KnowledgeBaseService struct {
	knowledgeBaseRepo *repository.KnowledgeBaseRepository
}

func NewKnowledgeBaseService() *KnowledgeBaseService {
	return &KnowledgeBaseService{
		knowledgeBaseRepo: repository.KnowledgeBaseRepo,
	}
}

type KnowledgeBaseListFilterArgs struct {
	IsPublic *bool `json:"is_public"`
}

type KnowledgeBaseListReq struct {
	CreatorID  *int64                       `json:"creator_id"`
	FilterArgs *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs   SortArgs                     `json:"sort_args"`
	Pagination Pagination                   `json:"pagination"`
}

func (s *KnowledgeBaseService) List(req *KnowledgeBaseListReq) ([]model.KnowledgeBase, error) {
	var knowledgeBases []model.KnowledgeBase

	return knowledgeBases, nil
}

type KnowledgeBaseListByExternalReq struct {
	CreatorExternalID string                       `json:"creator_external_id"`
	FilterArgs        *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                     `json:"sort_args"`
	Pagination        Pagination                   `json:"pagination"`
}

func (s *KnowledgeBaseService) ListByExternal(req *KnowledgeBaseListByExternalReq) ([]model.KnowledgeBase, error) {
	return nil, nil
}

var KnowledgeBaseSvc = NewKnowledgeBaseService()
