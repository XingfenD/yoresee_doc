package document_repo

import (
	"context"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentGetIDByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	tx         *gorm.DB
}

func (r *DocumentRepository) GetIDByExternalID(externalID string) *DocumentGetIDByExternalIDOperation {
	return &DocumentGetIDByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *DocumentGetIDByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentGetIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op DocumentGetIDByExternalIDOperation) query(tx *gorm.DB) (int64, error) {
	var id int64
	err := tx.Model(&model.Document{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (op *DocumentGetIDByExternalIDOperation) Exec(ctx context.Context) (int64, error) {
	if op.tx != nil {
		return op.query(op.tx)
	}

	extidKey := key.KeyIDByExternalID(key.KeyObjectTypeEnum_Doc, op.externalID)
	modelKey := key.KeyModelByExternalID(key.KeyObjectTypeEnum_Doc, op.externalID)

	id, err := cache_loader.NewCacheLoadOperation[int64](&op.repo.Loader).
		WithDefaultKeyAndParser(extidKey, cache_loader.ParseInt64).
		WithKeyAndParser(modelKey, cache_loader.ParseIDFromDocument).
		WithDBLoader(func() (*int64, error) {
			id, err := op.query(storage.DB)
			if err != nil {
				return nil, err
			}
			return &id, nil
		}).Exec(ctx)
	if err != nil {
		return 0, err
	}

	if id == nil {
		return 0, nil
	}

	return *id, nil
}
