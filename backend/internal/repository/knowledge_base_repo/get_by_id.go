package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type KnowledgeBaseGetByIDOperation struct {
	repo *KnowledgeBaseRepository
	id   int64
	tx   *gorm.DB
}

func (r *KnowledgeBaseRepository) GetByID(id int64) (op *KnowledgeBaseGetByIDOperation) {
	return &KnowledgeBaseGetByIDOperation{
		repo: KnowledgeBaseRepo,
		id:   id,
	}
}

func (op *KnowledgeBaseGetByIDOperation) WithTx(tx *gorm.DB) *KnowledgeBaseGetByIDOperation {
	op.tx = tx
	return op
}

func (op *KnowledgeBaseGetByIDOperation) Exec() (knowledgeBase *model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err = op.tx.First(knowledgeBase, "id = ?", op.id).Error
	return
}
