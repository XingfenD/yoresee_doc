package document_service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
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

func (s *DocumentService) SaveDocumentSnapshotAndContent(ctx context.Context, docExternalID string, state []byte, content string) (bool, error) {
	if len(state) == 0 {
		return false, status.StatusParamError
	}

	contentChanged := false
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		docModel, err := s.documentRepo.GetByExternalID(docExternalID).WithTx(tx).Exec(ctx)
		if err != nil {
			return status.StatusDocumentNotFound
		}

		if err := s.snapshotRepo.Save(docModel.ID, state).WithTx(tx).Exec(); err != nil {
			logrus.Errorf("[Service layer: DocumentService] SaveDocumentSnapshotAndContent save snapshot failed, doc_external_id=%s, doc_id=%d, err=%+v", docExternalID, docModel.ID, err)
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "save document snapshot failed")
		}

		if docModel.Content == content {
			return nil
		}

		versionModel := &model.DocumentVersion{
			DocumentID:    docModel.ID,
			Title:         docModel.Title,
			Content:       docModel.Content,
			UserID:        docModel.UserID,
			ChangeSummary: "Sync collab snapshot",
		}
		if err := s.docVersionRepo.Create(versionModel).WithTx(tx).Exec(); err != nil {
			logrus.Errorf("[Service layer: DocumentService] SaveDocumentSnapshotAndContent create version failed, doc_external_id=%s, doc_id=%d, err=%+v", docExternalID, docModel.ID, err)
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create document version failed")
		}

		nextDoc := &model.Document{
			ID:      docModel.ID,
			Content: content,
		}
		if err := s.documentRepo.Update(nextDoc).WithTx(tx).UpdateContent().Exec(); err != nil {
			logrus.Errorf("[Service layer: DocumentService] SaveDocumentSnapshotAndContent update content failed, doc_external_id=%s, doc_id=%d, err=%+v", docExternalID, docModel.ID, err)
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "update document content failed")
		}
		contentChanged = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return contentChanged, nil
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
