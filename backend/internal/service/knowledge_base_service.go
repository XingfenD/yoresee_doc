package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
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

func (s *KnowledgeBaseService) List(req *KnowledgeBaseListReq) ([]*dto.KnowledgeBaseResponse, error) {
	kbModels, _, err := s.knowledgeBaseRepo.List(&model.KnowledgeBase{}).
		WithCreatorID(req.CreatorID).
		WithIsPublic(req.FilterArgs.IsPublic).
		WithSort(req.SortArgs.Field, req.SortArgs.Desc).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()
	if err != nil {
		return nil, err
	}

	knowledgeBases := make([]*dto.KnowledgeBaseResponse, 0, len(kbModels))
	for _, kb := range kbModels {
		knowledgeBases = append(knowledgeBases, dto.NewKnowledgeBaseResponseFromModel(kb))
	}

	return knowledgeBases, nil
}

type KnowledgeBaseListByExternalReq struct {
	CreatorExternalID string                       `json:"creator_external_id"`
	FilterArgs        *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                     `json:"sort_args"`
	Pagination        Pagination                   `json:"pagination"`
}

func (s *KnowledgeBaseService) ListByExternal(req *KnowledgeBaseListByExternalReq) ([]*dto.KnowledgeBaseResponse, error) {
	var creatorID *int64
	if req.CreatorExternalID != "" {
		id, err := repository.UserRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
		if err != nil {
			return nil, status.StatusUserNotFound
		}
		creatorID = &id
	}

	return s.List(&KnowledgeBaseListReq{
		CreatorID:  creatorID,
		FilterArgs: req.FilterArgs,
		SortArgs:   req.SortArgs,
		Pagination: req.Pagination,
	})
}

var KnowledgeBaseSvc = NewKnowledgeBaseService()
