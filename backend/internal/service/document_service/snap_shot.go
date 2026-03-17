package document_service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"gorm.io/gorm"
)

func (s *DocumentService) SaveDocumentYjsSnapshot(ctx context.Context, docExternalID string, state []byte) error {
	if len(state) == 0 {
		return status.StatusParamError
	}
	docID, err := s.documentRepo.GetIDByExternalID(docExternalID).Exec(ctx)
	if err != nil {
		return status.StatusDocumentNotFound
	}

	return s.snapshotRepo.Save(docID, state).Exec()
}

func (s *DocumentService) GetDocumentYjsSnapshot(ctx context.Context, docExternalID string) ([]byte, error) {
	docID, err := s.documentRepo.GetIDByExternalID(docExternalID).Exec(ctx)
	if err != nil {
		return nil, status.StatusDocumentNotFound
	}

	snapshot, err := s.snapshotRepo.GetByDocID(docID).Exec()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return snapshot.YjsState, nil
}
