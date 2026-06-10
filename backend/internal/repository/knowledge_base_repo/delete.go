package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type DeleteKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Delete(knowledgeBase *model.KnowledgeBase) (op *DeleteKnowledgeBaseOperation) {
	return &DeleteKnowledgeBaseOperation{
		repo:          r,
		knowledgeBase: knowledgeBase,
	}
}

func (op *DeleteKnowledgeBaseOperation) WithTx(tx *gorm.DB) *DeleteKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *DeleteKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = op.repo.db
	}
	err := op.tx.Delete(op.knowledgeBase).Error
	return err
}
