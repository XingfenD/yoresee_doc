package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
)

func (s *DocumentService) GetDocumentByExternalID(ctx context.Context, externalID string) (*dto.DocumentResponse, error) {
	docModel, err := s.documentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return dto.NewDocumentResponseFromModel(docModel), nil
}
