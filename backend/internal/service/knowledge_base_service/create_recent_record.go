package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *KnowledgeBaseService) CreateRecentKnowledgeBase(req *dto.CreateRecentKnowledgeBaseRequest) error {
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	knowledgeBaseID, err := s.knowledgeBaseRepo.GetIDByExternalID(req.KnowledgeBaseExternalID).Exec()
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
