package document_service

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *DocumentService) RecordRecentDocument(req *dto.RecordRecentDocumentRequest) error {
	if req == nil || req.UserExternalID == "" || req.DocumentExternalID == "" {
		return status.StatusInternalParamsError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	docID, err := s.documentRepo.GetIDByExternalID(req.DocumentExternalID).Exec(context.Background())
	if err != nil || docID == 0 {
		return status.StatusDocumentNotFound
	}

	if err := s.documentRepo.UpsertRecentDocument(&model.RecentDocument{
		UserID:     userID,
		DocumentID: docID,
		AccessedAt: time.Now(),
	}).Exec(); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

func (s *DocumentService) ListRecentDocuments(req *dto.ListRecentDocumentsRequest) ([]*dto.DocumentResponse, int64, error) {
	if req == nil || req.UserExternalID == "" {
		return nil, 0, status.StatusInternalParamsError
	}
	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return nil, 0, status.StatusUserNotFound
	}

	recentRecords, total, err := s.documentRepo.ListRecentDocuments(userID).
		WithTimeRange(req.StartTime, req.EndTime).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		Exec()
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	if len(recentRecords) == 0 {
		return []*dto.DocumentResponse{}, total, nil
	}

	docIDs := make([]int64, 0, len(recentRecords))
	for _, record := range recentRecords {
		docIDs = append(docIDs, record.DocumentID)
	}

	docModels, err := s.documentRepo.MGetByIDs(docIDs)
	if err != nil {
		return nil, 0, status.StatusReadDBError
	}

	resp := make([]*dto.DocumentResponse, 0, len(docModels))
	for _, doc := range docModels {
		resp = append(resp, dto.NewDocumentResponseFromModel(doc))
	}

	return resp, total, nil
}
