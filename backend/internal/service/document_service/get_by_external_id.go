package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) GetDocumentByExternalID(ctx context.Context, externalID string) (*dto.DocumentResponse, error) {
	docModel, err := s.documentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] GetDocumentByExternalID failed, external_id=%s, err=%+v", externalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	return dto.NewDocumentResponseFromModel(docModel), nil
}
