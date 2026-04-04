package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *DocumentService) ListRecentTemplates(req *dto.ListRecentTemplatesRequest) ([]*dto.TemplateResponse, int64, error) {
	if req == nil || req.UserExternalID == "" {
		return nil, 0, status.StatusInternalParamsError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusUserNotFound
	}

	// TODO: filter the recent templates by kb
	recentRecords, total, err := s.templateRepo.ListRecentTemplates(userID).
		WithTimeRange(req.StartTime, req.EndTime).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	if len(recentRecords) == 0 {
		return []*dto.TemplateResponse{}, total, nil
	}

	tplIDs := make([]int64, 0, len(recentRecords))
	for _, r := range recentRecords {
		tplIDs = append(tplIDs, r.TemplateID)
	}

	tplModels, err := s.templateRepo.MGetByIDs(tplIDs).Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}
	tplMap := make(map[int64]*model.Template, len(tplModels))
	for _, tpl := range tplModels {
		tplMap[tpl.ID] = tpl
	}

	// map knowledge_base external ids
	kbExternalIDMap := make(map[int64]string)
	kbIDs := make([]int64, 0)
	for _, tpl := range tplModels {
		if tpl.KnowledgeBaseID == nil || *tpl.KnowledgeBaseID == 0 {
			continue
		}
		kbIDs = append(kbIDs, *tpl.KnowledgeBaseID)
	}
	if len(kbIDs) > 0 {
		kbs, err := s.kbRepo.MGetKnowledgeBaseByIDs(kbIDs).Exec()
		if err != nil {
			return nil, 0, status.StatusReadDBError
		}
		for _, kb := range kbs {
			kbExternalIDMap[kb.ID] = kb.ExternalID
		}
	}

	resp := make([]*dto.TemplateResponse, 0, len(recentRecords))
	for _, record := range recentRecords {
		tpl := tplMap[record.TemplateID]
		if tpl == nil {
			continue
		}
		kbExternalID := ""
		if tpl.KnowledgeBaseID != nil {
			kbExternalID = kbExternalIDMap[*tpl.KnowledgeBaseID]
		}
		resp = append(resp, &dto.TemplateResponse{
			ID:                      tpl.ID,
			Name:                    tpl.Name,
			Description:             tpl.Description,
			Content:                 tpl.Content,
			Type:                    doc_type_mapper.FromModelType(tpl.DocumentType),
			Scope:                   tpl.Scope,
			KnowledgeBaseExternalID: kbExternalID,
			Tags:                    tpl.Tags,
			CreatedAt:               tpl.CreatedAt,
			UpdatedAt:               tpl.UpdatedAt,
		})
	}

	return resp, total, nil
}
