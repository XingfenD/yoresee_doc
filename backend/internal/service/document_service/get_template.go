package document_service

import (
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"gorm.io/gorm"
)

func (s *DocumentService) GetTemplateByID(templateID int64) (*dto.TemplateResponse, error) {
	if templateID <= 0 {
		return nil, status.StatusParamError
	}

	tpl, err := s.templateRepo.GetByID(templateID).Exec()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.GenErrWithCustomMsg(status.StatusParamError, "template not found")
		}
		return nil, status.StatusReadDBError
	}
	if tpl == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "template not found")
	}

	kbExternalID := ""
	if tpl.KnowledgeBaseID != nil {
		kbs, err := s.kbRepo.MGetKnowledgeBaseByIDs([]int64{*tpl.KnowledgeBaseID}).Exec()
		if err != nil {
			return nil, status.StatusReadDBError
		}
		if len(kbs) > 0 {
			kbExternalID = kbs[0].ExternalID
		}
	}

	return &dto.TemplateResponse{
		ID:                      tpl.ID,
		Name:                    tpl.Name,
		Description:             tpl.Description,
		Content:                 tpl.Content,
		Scope:                   tpl.Scope,
		KnowledgeBaseExternalID: kbExternalID,
		Tags:                    tpl.Tags,
		CreatedAt:               tpl.CreatedAt,
		UpdatedAt:               tpl.UpdatedAt,
	}, nil
}
