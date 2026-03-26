package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

func (s *KnowledgeBaseService) ListByExternal(req *dto.KnowledgeBaseListByExternalReq) ([]*dto.KnowledgeBaseResponse, int64, error) {
	var creatorID *int64
	if req.CreatorExternalID != "" {
		id, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
		if err != nil {
			logrus.Errorf("[Service layer: KnowledgeBaseService] GetIDByExternalID failed, creator_external_id=%s, err=%+v", req.CreatorExternalID, err)
			return nil, 0, status.StatusUserNotFound
		}
		creatorID = &id
	}

	routerReq := buildKnowledgeBaseListReqFromExternal(req, creatorID)
	list, total, err := s.list(routerReq).WithDocumentExtend().WithUserExtend().ExecWithTotal()
	if err != nil {
		logrus.Errorf("[Service layer: KnowledgeBaseService] ListByExternal failed, err=%+v", err)
		return nil, 0, status.GenErrWithCustomMsg(err, "list knowledge bases failed")
	}
	return list, total, nil
}
