package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func buildKnowledgeBaseListReqFromExternal(req *dto.KnowledgeBaseListByExternalReq, creatorID *int64) *internal_dto.KnowledgeBaseListReq {
	if req == nil {
		return nil
	}

	return &internal_dto.KnowledgeBaseListReq{
		CreatorID:  creatorID,
		FilterArgs: req.FilterArgs,
		SortArgs:   req.SortArgs,
		Pagination: req.Pagination,
	}
}

func (s *KnowledgeBaseService) buildListKnowledgeBaseOperation(req *internal_dto.KnowledgeBaseListReq) (*knowledge_base_repo.ListKnowledgeBaseOperation, error) {
	if s == nil || s.knowledgeBaseRepo == nil {
		return nil, status.StatusServiceInternalError
	}
	if req == nil {
		return nil, status.StatusInternalParamsError
	}
	op := s.knowledgeBaseRepo.List(&model.KnowledgeBase{}).WithCreatorID(req.CreatorID)
	if req.FilterArgs != nil {
		op = op.WithIsPublic(req.FilterArgs.IsPublic).
			WithNameKeyword(req.FilterArgs.NameKeyword).
			WithCreateTimeRange(req.FilterArgs.CreateTimeRangeStart, req.FilterArgs.CreateTimeRangeEnd).
			WithUpdateTimeRange(req.FilterArgs.UpdateTimeRangeStart, req.FilterArgs.UpdateTimeRangeEnd)
	}
	op = op.WithSort(req.SortArgs.Field, req.SortArgs.Desc).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize)
	return op, nil
}
