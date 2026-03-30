package document_repo

import (
	"context"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DocumentGetByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	tx         *gorm.DB
}

func (r *DocumentRepository) GetByExternalID(externalID string) *DocumentGetByExternalIDOperation {
	return &DocumentGetByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *DocumentGetByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentGetByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetByExternalIDOperation) query(db *gorm.DB) (*model.Document, error) {
	var document model.Document
	var err error

	err = db.First(&document, "external_id = ?", op.externalID).Error

	return &document, err
}

func (op *DocumentGetByExternalIDOperation) Exec(ctx context.Context) (*model.Document, error) {
	var err error

	if op.tx != nil {
		return op.query(op.tx)
	}

	documentCacheKey := key.KeyModelByExternalID(key.KeyObjectTypeEnum_Doc, op.externalID)
	document, err := cache_loader.NewCacheLoadOperation[model.Document](&op.repo.Loader).
		WithDBLoader(func() (*model.Document, error) {
			return op.query(storage.DB)
		}).WithDefaultKeyAndParser(documentCacheKey, nil).
		Exec(ctx)

	if err != nil {
		logrus.Errorf("load data failed for DocumentGetByExternalIDOperation: %+v", err)
		return nil, err
	}

	return document, nil
}
