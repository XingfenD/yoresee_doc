package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
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

func (op *GetKnowledgeBaseByExternalIDOperation) Exec() (*model.KnowledgeBase, error) {
	var knowledgeBase model.KnowledgeBase
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.First(&knowledgeBase, "external_id = ?", op.externalID).Error
	if err != nil {
		return nil, err
	}
	return &knowledgeBase, nil
}
