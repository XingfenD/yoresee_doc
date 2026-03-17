package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type KnowledgeBaseGetIDByExternalIDOperation struct {
	repo       *KnowledgeBaseRepository
	externalID string
	tx         *gorm.DB
}

func (r *KnowledgeBaseRepository) GetIDByExternalID(externalID string) (op *KnowledgeBaseGetIDByExternalIDOperation) {
	return &KnowledgeBaseGetIDByExternalIDOperation{
		repo:       KnowledgeBaseRepo,
		externalID: externalID,
	}
}

func (op *KnowledgeBaseGetIDByExternalIDOperation) WithTx(tx *gorm.DB) *KnowledgeBaseGetIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *KnowledgeBaseGetIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.Model(&model.KnowledgeBase{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}
