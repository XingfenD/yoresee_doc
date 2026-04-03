package document_service

import (
	"context"
	"errors"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"gorm.io/gorm"
)

func (s *DocumentService) UpdateTemplateSettings(ctx context.Context, req *dto.UpdateTemplateSettingsRequest) error {
	if err := validateUpdateTemplateSettingsReq(req); err != nil {
		return status.GenErrWithCustomMsg(err, "invalid update template settings request")
	}
	_ = ctx

	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	tpl, err := s.templateRepo.GetByID(req.TemplateID).Exec()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.GenErrWithCustomMsg(status.StatusParamError, "template not found")
		}
		return status.StatusReadDBError
	}
	if tpl == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "template not found")
	}
	if tpl.UserID != userID {
		return status.StatusPermissionDenied
	}

	updateOp := s.templateRepo.Update(tpl)

	if req.Name != nil {
		nextName := strings.TrimSpace(*req.Name)
		if nextName == "" {
			return status.GenErrWithCustomMsg(status.StatusParamError, "template name is empty")
		}
		if len([]rune(nextName)) > 100 {
			nextName = truncateRunes(nextName, 100)
		}
		tpl.Name = nextName
		updateOp.UpdateName()
	}

	if req.Description != nil {
		tpl.Description = strings.TrimSpace(*req.Description)
		updateOp.UpdateDescription()
	}

	if req.IsPublic != nil {
		tpl.IsPublic = *req.IsPublic
		if *req.IsPublic {
			tpl.Scope = "system"
		} else if tpl.KnowledgeBaseID != nil {
			tpl.Scope = "knowledge_base"
		} else {
			tpl.Scope = "private"
		}
		updateOp.UpdateScope().UpdateIsPublic()
	}

	if err := updateOp.Exec(); err != nil {
		return status.StatusWriteDBError
	}

	return nil
}
