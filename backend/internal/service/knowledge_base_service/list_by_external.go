package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *KnowledgeBaseService) ListByExternal(req *dto.KnowledgeBaseListByExternalReq) ([]*dto.KnowledgeBaseResponse, int64, error) {
	var creatorID *int64
	if req.CreatorExternalID != "" {
		id, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
		if err != nil {
			return nil, 0, status.StatusUserNotFound
		}
		creatorID = &id
	}

	routerReq := buildKnowledgeBaseListReqFromExternal(req, creatorID)
	return s.list(routerReq).WithDocumentExtend().WithUserExtend().ExecWithTotal()
}
