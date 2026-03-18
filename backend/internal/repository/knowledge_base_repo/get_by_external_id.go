package knowledge_base_repo

import (
	"context"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GetKnowledgeBaseByExternalIDOperation struct {
	repo       *KnowledgeBaseRepository
	externalID string
	tx         *gorm.DB
}

func (r *KnowledgeBaseRepository) GetByExternalID(externalID string) (op *GetKnowledgeBaseByExternalIDOperation) {
	return &GetKnowledgeBaseByExternalIDOperation{
		repo:       KnowledgeBaseRepo,
		externalID: externalID,
	}
}

func (op *GetKnowledgeBaseByExternalIDOperation) WithTx(tx *gorm.DB) *GetKnowledgeBaseByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *GetKnowledgeBaseByExternalIDOperation) query(db *gorm.DB) (*model.KnowledgeBase, error) {
	var knowledgeBase model.KnowledgeBase
	err := db.First(&knowledgeBase, "external_id = ?", op.externalID).Error
	return &knowledgeBase, err
}

func (op *GetKnowledgeBaseByExternalIDOperation) Exec() (*model.KnowledgeBase, error) {
	if op.tx != nil {
		return op.query(op.tx)
	}

	knowledgeBaseCacheKey := cache.KeyModelByExternalID(cache.KeyObjectTypeEnum_KnowledgeBase, op.externalID)
	knowledgeBase, err := cache_loader.NewCacheLoadOperation[model.KnowledgeBase](&op.repo.Loader).
		WithDBLoader(func() (*model.KnowledgeBase, error) {
			return op.query(storage.DB)
		}).WithDefaultKeyAndParser(knowledgeBaseCacheKey, nil).
		Exec(context.Background())

	if err != nil {
		logrus.Errorf("load data failed for GetKnowledgeBaseByExternalIDOperation: %+v", err)
		return nil, err
	}

	return knowledgeBase, nil
}
