package document_repo

import (
	"context"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	Loader cache_loader.Loader
}

var DocumentRepo DocumentRepository

func Init(redis *redis.Client) {
	DocumentRepo.Loader = *cache_loader.NewLoader(redis)
}

type DocumentDeleteOperation struct {
	repo *DocumentRepository
	id   int64
	tx   *gorm.DB
	ctx  context.Context
}

func (r *DocumentRepository) Delete(id int64) *DocumentDeleteOperation {
	return &DocumentDeleteOperation{
		repo: r,
		id:   id,
		ctx:  context.Background(),
	}
}

func (op *DocumentDeleteOperation) WithTx(tx *gorm.DB) *DocumentDeleteOperation {
	op.tx = tx
	return op
}

func (op *DocumentDeleteOperation) WithContext(ctx context.Context) *DocumentDeleteOperation {
	op.ctx = ctx
	return op
}

func (op *DocumentDeleteOperation) Exec() error {
	var doc model.Document
	var db *gorm.DB
	var err error

	if op.tx != nil {
		db = op.tx
	} else {
		db = storage.DB
	}

	if err := db.First(&doc, op.id).Error; err != nil {
		return err
	}

	if op.tx == nil {
		db = db.Begin()
		defer func() {
			if err != nil {
				db.Rollback()
			} else {
				db.Commit()
			}
		}()
	}

	if err := db.Delete(&model.Document{}, op.id).Error; err != nil {
		return err
	}

	if op.tx == nil {
		docCacheKey := key.KeyModelByExternalID(key.KeyObjectTypeEnum_Doc, doc.ExternalID)
		if err := storage.KVS.Del(op.ctx, docCacheKey).Err(); err != nil {
		}

		if err := op.repo.BumpSubtreeVersionsByPath(op.ctx, doc.Path); err != nil {
		}
	}

	return nil
}
