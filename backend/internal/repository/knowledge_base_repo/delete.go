package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DeleteKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Delete(knowledgeBase *model.KnowledgeBase) (op *DeleteKnowledgeBaseOperation) {
	return &DeleteKnowledgeBaseOperation{
		repo:          KnowledgeBaseRepo,
		knowledgeBase: knowledgeBase,
	}
}

func (op *DeleteKnowledgeBaseOperation) WithTx(tx *gorm.DB) *DeleteKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *DeleteKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.Delete(op.knowledgeBase).Error
	return err
}
