package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/repository/knowledge_base_repo"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

type KnowledgeBaseService struct {
	knowledgeBaseRepo *knowledge_base_repo.KnowledgeBaseRepository
	userRepo          *user_repo.UserRepository
}

func NewKnowledgeBaseService() *KnowledgeBaseService {
	return &KnowledgeBaseService{
		knowledgeBaseRepo: knowledge_base_repo.KnowledgeBaseRepo,
		userRepo:          user_repo.UserRepo,
	}
}

func (s *KnowledgeBaseService) GetIDByExternalID(externalID string) (int64, error) {
	id, err := s.knowledgeBaseRepo.GetIDByExternalID(externalID).Exec()
	if err != nil {
		return 0, status.StatusKnowledgeBaseNotFound
	}
	return id, nil
}

var KnowledgeBaseSvc = NewKnowledgeBaseService()
