package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *DocumentService) UpdateDocumentSettings(ctx context.Context, req *dto.UpdateDocumentSettingsRequest) (*dto.DocumentSettingsResponse, error) {
	if req == nil || req.ExternalID == "" {
		return nil, status.StatusParamError
	}

	cacheKey := key.KeyModelByExternalID(key.KeyObjectTypeEnum_Doc, req.ExternalID)
	err := cache.DoubleDelete(
		context.Background(),
		func() error {
			return utils.WithTransaction(func(tx *gorm.DB) error {
				oldDoc, err := s.documentRepo.GetByExternalID(req.ExternalID).WithTx(tx).Exec(ctx)
				if err != nil {
					return status.StatusDocumentNotFound
				}

				docModel := &model.Document{
					ID:       oldDoc.ID,
					IsPublic: req.IsPublic,
				}

				if err := s.documentRepo.Update(docModel).WithTx(tx).UpdateIsPublic().Exec(); err != nil {
					logrus.Errorf("[Service layer: DocumentService] update document settings failed, external_id=%s, err=%+v", req.ExternalID, err)
					return status.GenErrWithCustomMsg(status.StatusWriteDBError, "update document settings failed")
				}

				return nil
			})
		},
		cacheKey,
	)
	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] UpdateDocumentSettings failed, external_id=%s, err=%+v", req.ExternalID, err)
		return nil, status.GenErrWithCustomMsg(err, "update document settings failed")
	}

	if err := domain_event.PublishDocumentUpsertEvent(ctx, req.ExternalID); err != nil {
		logrus.Warnf("[Service layer: DocumentService] publish search sync event failed, external_id=%s, err=%+v", req.ExternalID, err)
	}

	return &dto.DocumentSettingsResponse{
		IsPublic: req.IsPublic,
	}, nil
}
