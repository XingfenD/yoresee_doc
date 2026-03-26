package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

type KnowledgeBaseGetByExternalIDOperation struct {
	withUserExtend     bool
	withDocumentExtend bool
	req                *dto.KnowledgeBaseGetByExternalIDReq
	srvc               *KnowledgeBaseService
}

func (s *KnowledgeBaseService) GetByExternalID(req *dto.KnowledgeBaseGetByExternalIDReq) *KnowledgeBaseGetByExternalIDOperation {
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
		logrus.Errorf("[Service layer: KnowledgeBaseService] GetByExternalID failed, external_id=%s, err=%+v", op.req.KnowledgeBaseExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusKnowledgeBaseNotFound, "knowledge base not found")
	}
	extendDTO := &dto.KnowledgeBaseExtend{}

	if op.withUserExtend {
		userModel, err := op.srvc.userRepo.GetByID(kbModel.CreatorUserID).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: KnowledgeBaseService] Get creator user failed, user_id=%d, err=%+v", kbModel.CreatorUserID, err)
			return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "query creator failed")
		}
		extendDTO.CreatorUserExternalID = userModel.ExternalID
		extendDTO.CreatorName = userModel.Username
	}

	if op.withDocumentExtend {
		ids := []int64{
			kbModel.ID,
		}
		count, err := op.srvc.knowledgeBaseRepo.MGetKnowledgeBaseDocumentsCount(ids).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: KnowledgeBaseService] MGetKnowledgeBaseDocumentsCount failed, kb_id=%d, err=%+v", kbModel.ID, err)
			return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "query document count failed")
		}
		extendDTO.DocumentsCount = count[kbModel.ID]
	}

	return dto.NewKnowledgeBaseResponseFromModel(
		kbModel, extendDTO,
	), nil
}
