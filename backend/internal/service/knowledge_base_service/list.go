package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/bytedance/gg/gslice"
	"github.com/sirupsen/logrus"
)

type knowledgeBaseListOperation struct {
	req  *internal_dto.KnowledgeBaseListReq
	srvc *KnowledgeBaseService

	withDocumentExtend bool
	withUserExtend     bool
}

func (s *KnowledgeBaseService) list(req *internal_dto.KnowledgeBaseListReq) *knowledgeBaseListOperation {
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
			return status.GenErrWithCustomMsg(status.StatusReadDBError, "query knowledge base document count failed")
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
		usersMapByUserID, err := op.srvc.userRepo.MGetUserByID(uniqUserID).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: knowledgeBaseList]: MGetUserByID failed, err: %+v", err)
			return status.GenErrWithCustomMsg(status.StatusReadDBError, "query creator users failed")
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
		logrus.Errorf("[Service layer: knowledgeBaseList]: buildListKnowledgeBaseOperation failed, err: %+v", err)
		return nil, status.GenErrWithCustomMsg(err, "build knowledge base list operation failed")
	}
	kbModels, err := listOp.Exec()
	if err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: list operation exec failed, err: %+v", err)
		return nil, status.GenErrWithCustomMsg(err, "list knowledge bases failed")
	}

	kbExtendMapByID := make(map[int64]*dto.KnowledgeBaseExtend)

	if err := op.documentExtend(kbModels, kbExtendMapByID); err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: documentExtend failed, err: %+v", err)
		return nil, status.GenErrWithCustomMsg(err, "build knowledge base document extension failed")
	}
	if err := op.userExtend(kbModels, kbExtendMapByID); err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: userExtend failed, err: %+v", err)
		return nil, status.GenErrWithCustomMsg(err, "build knowledge base user extension failed")
	}

	knowledgeBases := make([]*dto.KnowledgeBaseResponse, 0, len(kbModels))
	for _, kb := range kbModels {
		knowledgeBases = append(knowledgeBases, dto.NewKnowledgeBaseResponseFromModel(kb, kbExtendMapByID[kb.ID]))
	}

	return knowledgeBases, nil
}

func (op *knowledgeBaseListOperation) ExecWithTotal() ([]*dto.KnowledgeBaseResponse, int64, error) {
	listOp, err := op.srvc.buildListKnowledgeBaseOperation(op.req)
	if err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: buildListKnowledgeBaseOperation failed, err: %+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build knowledge base list operation failed")
	}
	kbModels, total, err := listOp.ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: list operation exec with total failed, err: %+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "list knowledge bases failed")
	}

	kbExtendMapByID := make(map[int64]*dto.KnowledgeBaseExtend)

	if err := op.documentExtend(kbModels, kbExtendMapByID); err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: documentExtend failed, err: %+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build knowledge base document extension failed")
	}
	if err := op.userExtend(kbModels, kbExtendMapByID); err != nil {
		logrus.Errorf("[Service layer: knowledgeBaseList]: userExtend failed, err: %+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "build knowledge base user extension failed")
	}

	knowledgeBases := make([]*dto.KnowledgeBaseResponse, 0, len(kbModels))
	for _, kb := range kbModels {
		knowledgeBases = append(knowledgeBases, dto.NewKnowledgeBaseResponseFromModel(kb, kbExtendMapByID[kb.ID]))
	}

	return knowledgeBases, total, nil
}
