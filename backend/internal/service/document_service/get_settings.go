package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func (s *DocumentService) GetDocumentSettings(ctx context.Context, req *dto.GetDocumentSettingsRequest) (*dto.DocumentSettingsResponse, error) {
	if req == nil || req.ExternalID == "" {
		return nil, status.StatusParamError
	}

	// Settings page expects real-time value after save/update, so read from DB directly
	// to avoid returning a stale cached model.
	docModel, err := s.documentRepo.GetByExternalID(req.ExternalID).WithTx(storage.DB).Exec(ctx)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] GetDocumentSettings failed, external_id=%s, err=%+v", req.ExternalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusDocumentNotFound, "document not found")
	}

	// TODO: merge fields from document_setting table when that table is introduced.
	return &dto.DocumentSettingsResponse{
		IsPublic: docModel.IsPublic,
	}, nil
}
