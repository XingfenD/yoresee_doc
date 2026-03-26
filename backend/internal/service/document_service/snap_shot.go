package document_service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
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

	if err := s.snapshotRepo.Save(docID, state).Exec(); err != nil {
		logrus.Errorf("[Service layer: DocumentService] SaveDocumentYjsSnapshot failed, doc_external_id=%s, doc_id=%d, err=%+v", docExternalID, docID, err)
		return status.GenErrWithCustomMsg(status.StatusWriteDBError, "save document snapshot failed")
	}
	return nil
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
		logrus.Errorf("[Service layer: DocumentService] GetDocumentYjsSnapshot failed, doc_external_id=%s, doc_id=%d, err=%+v", docExternalID, docID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusReadDBError, "get document snapshot failed")
	}
	return snapshot.YjsState, nil
}
