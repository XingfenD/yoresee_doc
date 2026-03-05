package service

import (
	"github.com/bytedance/gg/gslice"
	"github.com/sirupsen/logrus"

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

type knowledgeBaseListReq struct {
	CreatorID  *int64                       `json:"creator_id"`
	FilterArgs *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs   SortArgs                     `json:"sort_args"`
	Pagination Pagination                   `json:"pagination"`
}

func (s *KnowledgeBaseService) buildListKnowledgeBaseOperation(req *knowledgeBaseListReq) (*repository.ListKnowledgeBaseOperation, error) {
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

type knowledgeBaseListOperation struct {
	req  *knowledgeBaseListReq
	srvc *KnowledgeBaseService

	withDocumentExtend bool
	withUserExtend     bool
}

func (s *KnowledgeBaseService) list(req *knowledgeBaseListReq) *knowledgeBaseListOperation {
	return &knowledgeBaseListOperation{
		req:  req,
		srvc: s,
	}
}

func (op *knowledgeBaseListOperation) WithDocumentExtend() *knowledgeBaseListOperation {
	op.withDocumentExtend = true
	return op
}

func (op *knowledgeBaseListOperation) WithUserExtend() *knowledgeBaseListOperation {
	op.withUserExtend = true
	return op
}

func (op *knowledgeBaseListOperation) documentExtend(kbModels []*model.KnowledgeBase, kbExtendMapByID map[int64]*dto.KnowledgeBaseExtend) error {
	if op.withDocumentExtend {
		kbIDs := gslice.Map(kbModels, func(kbModel *model.KnowledgeBase) int64 {
			return kbModel.ID
		})
		countMapByKbID, err := op.srvc.knowledgeBaseRepo.MGetKnowledgeBaseDocumentsCount(kbIDs).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: knowledgeBaseList]: MGetKnowledgeBaseDocumentsCount failed, err: %+v", err)
			return err
		}
		for id, count := range countMapByKbID {
			if _, ok := kbExtendMapByID[id]; !ok {
				kbExtendMapByID[id] = &dto.KnowledgeBaseExtend{}
			}
			kbExtendMapByID[id].DocumentsCount = count
		}
	}
	return nil
}

func (op *knowledgeBaseListOperation) userExtend(kbModels []*model.KnowledgeBase, kbExtendMapByID map[int64]*dto.KnowledgeBaseExtend) error {
	if op.withUserExtend {
		// collect all user_id
		allUserID := gslice.Map(kbModels, func(kb *model.KnowledgeBase) int64 {
			return kb.CreatorUserID
		})
		uniqUserID := gslice.Uniq(allUserID)
		usersMapByUserID, err := op.srvc.userSrvc.userRepo.MGetUserByID(uniqUserID).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: knowledgeBaseList]: MGetUserByID failed, err: %+v", err)
			return err
		}
		for _, kbModel := range kbModels {
			if _, ok := kbExtendMapByID[kbModel.ID]; !ok {
				kbExtendMapByID[kbModel.ID] = &dto.KnowledgeBaseExtend{}
			}
			if user, ok := usersMapByUserID[kbModel.CreatorUserID]; ok {
				kbExtendMapByID[kbModel.ID].CreatorUserExternalID = user.ExternalID
				kbExtendMapByID[kbModel.ID].CreatorName = user.Username
			}
		}
	}
	return nil
}

func (op *knowledgeBaseListOperation) Exec() ([]*dto.KnowledgeBaseResponse, error) {
	listOp, err := op.srvc.buildListKnowledgeBaseOperation(op.req)
	if err != nil {
		return nil, err
	}
	kbModels, err := listOp.Exec()
	if err != nil {
		return nil, err
	}

	kbExtendMapByID := make(map[int64]*dto.KnowledgeBaseExtend)

	op.documentExtend(kbModels, kbExtendMapByID)
	op.userExtend(kbModels, kbExtendMapByID)

	knowledgeBases := make([]*dto.KnowledgeBaseResponse, 0, len(kbModels))
	for _, kb := range kbModels {
		knowledgeBases = append(knowledgeBases, dto.NewKnowledgeBaseResponseFromModel(kb, kbExtendMapByID[kb.ID]))
	}

	return knowledgeBases, nil
}

func (op *knowledgeBaseListOperation) ExecWithTotal() ([]*dto.KnowledgeBaseResponse, int64, error) {
	listOp, err := op.srvc.buildListKnowledgeBaseOperation(op.req)
	if err != nil {
		return nil, 0, err
	}
	kbModels, total, err := listOp.ExecWithTotal()
	if err != nil {
		return nil, 0, err
	}

	kbExtendMapByID := make(map[int64]*dto.KnowledgeBaseExtend)

	op.documentExtend(kbModels, kbExtendMapByID)
	op.userExtend(kbModels, kbExtendMapByID)

	knowledgeBases := make([]*dto.KnowledgeBaseResponse, 0, len(kbModels))
	for _, kb := range kbModels {
		knowledgeBases = append(knowledgeBases, dto.NewKnowledgeBaseResponseFromModel(kb, kbExtendMapByID[kb.ID]))
	}

	return knowledgeBases, total, nil
}

type KnowledgeBaseListByExternalReq struct {
	CreatorExternalID string                       `json:"creator_external_id"`
	FilterArgs        *KnowledgeBaseListFilterArgs `json:"filter_args"`
	SortArgs          SortArgs                     `json:"sort_args"`
	Pagination        Pagination                   `json:"pagination"`
}

func buildKnowledgeBaseListReqFromExternal(req *KnowledgeBaseListByExternalReq, creatorID *int64) *knowledgeBaseListReq {
	if req == nil {
		return nil
	}

	return &knowledgeBaseListReq{
		CreatorID:  creatorID,
		FilterArgs: req.FilterArgs,
		SortArgs:   req.SortArgs,
		Pagination: req.Pagination,
	}
}

func (s *KnowledgeBaseService) ListByExternal(req *KnowledgeBaseListByExternalReq) ([]*dto.KnowledgeBaseResponse, int64, error) {
	var creatorID *int64
	if req.CreatorExternalID != "" {
		id, err := repository.UserRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
		if err != nil {
			return nil, 0, status.StatusUserNotFound
		}
		creatorID = &id
	}

	routerReq := buildKnowledgeBaseListReqFromExternal(req, creatorID)
	return s.list(routerReq).WithDocumentExtend().WithUserExtend().ExecWithTotal()
}

func (s *KnowledgeBaseService) GetIDByExternalID(externalID string) (int64, error) {
	id, err := s.knowledgeBaseRepo.GetIDByExternalID(externalID).Exec()
	if err != nil {
		return 0, status.StatusKnowledgeBaseNotFound
	}
	return id, nil
}

type KnowledgeBaseGetByExternalIDReq struct {
	KnowledgeBaseExternalID string `json:"knowledge_base_external_id"`
}

type KnowledgeBaseGetByExternalIDOperation struct {
	withUserExtend     bool
	withDocumentExtend bool
	req                *KnowledgeBaseGetByExternalIDReq
	srvc               *KnowledgeBaseService
}

func (s *KnowledgeBaseService) GetByExternalID(req *KnowledgeBaseGetByExternalIDReq) *KnowledgeBaseGetByExternalIDOperation {
	return &KnowledgeBaseGetByExternalIDOperation{
		req:  req,
		srvc: s,
	}
}

func (op *KnowledgeBaseGetByExternalIDOperation) WithExtend() *KnowledgeBaseGetByExternalIDOperation {
	op.withUserExtend = true
	op.withDocumentExtend = true
	return op
}

func (op *KnowledgeBaseGetByExternalIDOperation) WithUserExtend() *KnowledgeBaseGetByExternalIDOperation {
	op.withUserExtend = true
	return op
}

func (op *KnowledgeBaseGetByExternalIDOperation) WithDocumentExtend() *KnowledgeBaseGetByExternalIDOperation {
	op.withDocumentExtend = true
	return op
}

func (op *KnowledgeBaseGetByExternalIDOperation) Exec() (*dto.KnowledgeBaseResponse, error) {
	kbModel, err := op.srvc.knowledgeBaseRepo.GetByExternalID(op.req.KnowledgeBaseExternalID).Exec()
	if err != nil {
		return nil, err
	}
	extendDTO := &dto.KnowledgeBaseExtend{}

	if op.withUserExtend {
		userModel, err := op.srvc.userSrvc.GetByID(kbModel.CreatorUserID)
		if err != nil {
			return nil, err
		}
		extendDTO.CreatorUserExternalID = userModel.ExternalID
		extendDTO.CreatorName = userModel.Username
	}

	if op.withDocumentExtend {
		count, err := op.srvc.knowledgeBaseRepo.GetKnowledgeBaseDocumentsCount(kbModel.ID).Exec()
		if err != nil {
			return nil, err
		}
		extendDTO.DocumentsCount = count
	}

	return dto.NewKnowledgeBaseResponseFromModel(
		kbModel, extendDTO,
	), nil
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
