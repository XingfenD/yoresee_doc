package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type CreateKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Create(knowledgeBase *model.KnowledgeBase) (op *CreateKnowledgeBaseOperation) {
	return &CreateKnowledgeBaseOperation{
		repo:          r,
		knowledgeBase: knowledgeBase,
	}
}

func (op *CreateKnowledgeBaseOperation) WithTx(tx *gorm.DB) *CreateKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *CreateKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = op.repo.db
	}

	return op.tx.Create(op.knowledgeBase).Error
}
