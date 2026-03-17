package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CreateKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Create(knowledgeBase *model.KnowledgeBase) (op *CreateKnowledgeBaseOperation) {
	return &CreateKnowledgeBaseOperation{
		repo:          KnowledgeBaseRepo,
		knowledgeBase: knowledgeBase,
	}
}

func (op *CreateKnowledgeBaseOperation) WithTx(tx *gorm.DB) *CreateKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *CreateKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Create(op.knowledgeBase).Error
}
