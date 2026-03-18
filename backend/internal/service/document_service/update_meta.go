package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"gorm.io/gorm"
)

func (s *DocumentService) UpdateDocumentMeta(ctx context.Context, req *dto.UpdateDocumentMetaRequest) (bool, error) {
	if err := validateUpdateDocumentMetaReq(req); err != nil {
		return false, err
	}

	cacheKey := cache.KeyModelByExternalID(cache.KeyObjectTypeEnum_Doc, req.ExternalID)
	err := cache.DoubleDelete(
		ctx,
		func() error {
			return utils.WithTransaction(func(tx *gorm.DB) error {
				oldDoc, err := s.documentRepo.GetByExternalID(req.ExternalID).WithTx(tx).Exec(ctx)
				if err != nil {
					return status.StatusDocumentNotFound
				}

				docModel := &model.Document{ID: oldDoc.ID}
				op := s.documentRepo.Update(docModel).WithTx(tx)

				if req.Title != nil {
					docModel.Title = *req.Title
					op.UpdateTitle()
				}
				if req.Summary != nil {
					docModel.Summary = *req.Summary
					op.UpdateSummary()
				}
				if req.Tags != nil {
					docModel.Tags = *req.Tags
					op.UpdateTags()
				}
				if req.Status != nil {
					docModel.Status = *req.Status
					op.UpdateStatus()
				}

				if err := op.Exec(); err != nil {
					return err
				}

				versionModel := &model.DocumentVersion{
					DocumentID:    oldDoc.ID,
					Title:         oldDoc.Title,
					Content:       oldDoc.Content,
					UserID:        oldDoc.UserID,
					ChangeSummary: "Update document meta",
				}
				if err := s.docVersionRepo.Create(versionModel).WithTx(tx).Exec(); err != nil {
					return err
				}

				return nil
			})
		},
		cacheKey,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
