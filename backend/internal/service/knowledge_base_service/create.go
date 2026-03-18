package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
)

func (s *KnowledgeBaseService) Create(req *dto.CreateKnowledgeBaseRequest) (*dto.CreateKnowledgeBaseResponse, error) {
	if req == nil || req.Name == "" {
		return nil, status.StatusParamError
	}

	userID, err := s.userRepo.GetIDByExternalID(req.CreatorExternalID).Exec()
	if err != nil {
		return nil, status.StatusUserNotFound
	}

	kbModel := &model.KnowledgeBase{
		ExternalID:    utils.GenerateExternalID("kb"),
		Name:          req.Name,
		Description:   req.Description,
		Cover:         req.Cover,
		IsPublic:      req.IsPublic,
		CreatorUserID: userID,
	}

	if err := s.knowledgeBaseRepo.Create(kbModel).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}

	return &dto.CreateKnowledgeBaseResponse{ExternalID: kbModel.ExternalID}, nil
}
