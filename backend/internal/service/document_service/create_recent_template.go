package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *DocumentService) CreateRecentTemplate(req *dto.CreateRecentTemplateRequest) error {
	if req == nil || req.UserExternalID == "" || req.TemplateID <= 0 {
		return status.StatusInternalParamsError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	err = s.templateRepo.CreateRecentTemplate(&model.RecentTemplate{
		UserID:     userID,
		TemplateID: req.TemplateID,
		AccessedAt: req.AccessTime,
	}).Exec()
	if err != nil {
		return status.StatusWriteDBError
	}
	return nil
}
