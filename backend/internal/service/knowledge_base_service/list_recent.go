package knowledge_base_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *KnowledgeBaseService) ListRecentKnowledgeBases(req *dto.ListRecentKnowledgeBasesRequest) ([]*dto.KnowledgeBaseResponse, int64, error) {
	if req == nil || req.UserExternalID == "" {
		return nil, 0, status.StatusInternalParamsError
	}

	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusUserNotFound
	}

	recentRecords, total, err := s.knowledgeBaseRepo.ListRecentKnowledgeBases(userID).
		WithTimeRange(req.StartTime, req.EndTime).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	if len(recentRecords) == 0 {
		return []*dto.KnowledgeBaseResponse{}, total, nil
	}

	kbIDs := make([]int64, 0, len(recentRecords))
	for _, r := range recentRecords {
		kbIDs = append(kbIDs, r.KnowledgeBaseID)
	}

	kbModels, err := s.knowledgeBaseRepo.MGetKnowledgeBaseByIDs(kbIDs).Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	kbMap := make(map[int64]*model.KnowledgeBase, len(kbModels))
	for _, kb := range kbModels {
		kbMap[kb.ID] = kb
	}

	kbExtendMapByID := make(map[int64]*dto.KnowledgeBaseExtend)
	op := &knowledgeBaseListOperation{
		srvc:               s,
		withDocumentExtend: true,
		withUserExtend:     true,
	}
	_ = op.documentExtend(kbModels, kbExtendMapByID)
	_ = op.userExtend(kbModels, kbExtendMapByID)

	resp := make([]*dto.KnowledgeBaseResponse, 0, len(recentRecords))
	for _, record := range recentRecords {
		kb := kbMap[record.KnowledgeBaseID]
		if kb == nil {
			continue
		}
		resp = append(resp, dto.NewKnowledgeBaseResponseFromModel(kb, kbExtendMapByID[kb.ID]))
	}

	return resp, total, nil
}
