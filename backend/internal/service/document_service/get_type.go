package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *DocumentService) GetDocumentTypeByExternalID(ctx context.Context, externalID string) (model.DocumentType, error) {
	docModel, err := s.documentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		return "", status.StatusDocumentNotFound
	}
	return docModel.Type, nil
}
