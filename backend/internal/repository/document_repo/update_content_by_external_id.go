package document_repo

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentUpdateContentByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	content    string
	tx         *gorm.DB
}

func (r *DocumentRepository) UpdateContentByExternalID(externalID, content string) *DocumentUpdateContentByExternalIDOperation {
	return &DocumentUpdateContentByExternalIDOperation{
		repo:       r,
		externalID: externalID,
		content:    content,
	}
}

func (op *DocumentUpdateContentByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentUpdateContentByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentUpdateContentByExternalIDOperation) Exec(ctx context.Context) error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	docModelCacheKey := cache.KeyModelByExternalID(cache.KeyObjectTypeEnum_Doc, op.externalID)
	return cache.DoubleDelete(
		ctx,
		func() error {
			return op.tx.Model(&model.Document{}).
				Where("external_id = ?", op.externalID).
				Select("content").
				Updates(map[string]interface{}{"content": op.content}).Error
		},
		docModelCacheKey,
	)

}
