package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

type KnowledgeBaseService struct {
	knowledgeBaseRepo *repository.KnowledgeBaseRepository
	userSrvc          *UserService
}

func NewKnowledgeBaseService() *KnowledgeBaseService {
	return &KnowledgeBaseService{
		knowledgeBaseRepo: repository.KnowledgeBaseRepo,
		userSrvc:          UserSvc,
	}
}

type KnowledgeBaseListFilterArgs struct {
	IsPublic             *bool   `json:"is_public"`
	NameKeyword          *string `json:"name_keyword"`
	CreateTimeRangeStart *string `json:"create_time_range_start"`
	CreateTimeRangeEnd   *string `json:"create_time_range_end"`
	UpdateTimeRangeStart *string `json:"update_time_range_start"`
	UpdateTimeRangeEnd   *string `json:"update_time_range_end"`
}

type KnowledgeBaseListReq struct {
	CreatorID  *int64                       `json:"creator_id"`
	FilterArgs *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs   SortArgs                     `json:"sort_args"`
	Pagination Pagination                   `json:"pagination"`
}

func (s *KnowledgeBaseService) buildListKnowledgeBaseOperation(req *KnowledgeBaseListReq) (*repository.ListKnowledgeBaseOperation, error) {
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

func (s *KnowledgeBaseService) List(req *KnowledgeBaseListReq) ([]*dto.KnowledgeBaseResponse, error) {
	listOp, err := s.buildListKnowledgeBaseOperation(req)
	if err != nil {
		return nil, err
	}
	kbModels, _, err := listOp.ExecWithTotal()
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

func (s *KnowledgeBaseService) GetIDByExternalID(externalID string) (int64, error) {
	id, err := s.knowledgeBaseRepo.GetIDByExternalID(externalID).Exec()
	if err != nil {
		return 0, status.StatusKnowledgeBaseNotFound
	}
	return id, nil
}

type KnowledgeBaseGetReq struct {
	KnowledgeBaseExternalID string `json:"knowledge_base_external_id"`
}

func (s *KnowledgeBaseService) Get(req *KnowledgeBaseGetReq) (*dto.KnowledgeBaseResponse, error) {
	kbModel, err := s.knowledgeBaseRepo.GetByExternalID(req.KnowledgeBaseExternalID).Exec()
	if err != nil {
		return nil, err
	}

	return dto.NewKnowledgeBaseResponseFromModel(kbModel), nil
}

func (s *KnowledgeBaseService) CreateRecentKnowledgeBase(req *dto.CreateRecentKnowledgeBaseRequest) error {
	userID, err := s.userSrvc.GetIDByExternalID(req.UserExternalID)
	if err != nil {
		return status.StatusUserNotFound
	}

	knowledgeBaseID, err := s.GetIDByExternalID(req.KnowledgeBaseExternalID)
	if err != nil {
		return status.StatusKnowledgeBaseNotFound
	}

	err = s.knowledgeBaseRepo.CreateRecentKnowledgeBase(&model.RecentKnowledgeBase{
		UserID:          userID,
		KnowledgeBaseID: knowledgeBaseID,
		AccessedAt:      req.AssessTime,
	}).Exec()

	if err != nil {
		return status.StatusWriteDBError
	}

	return nil
}

var KnowledgeBaseSvc = NewKnowledgeBaseService()
