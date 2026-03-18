package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type MGetKnowledgeBaseByIDsOperation struct {
	repo *KnowledgeBaseRepository
	ids  []int64
	tx   *gorm.DB
}

func (r *KnowledgeBaseRepository) MGetKnowledgeBaseByIDs(ids []int64) *MGetKnowledgeBaseByIDsOperation {
	return &MGetKnowledgeBaseByIDsOperation{
		repo: r,
		ids:  ids,
	}
}

func (op *MGetKnowledgeBaseByIDsOperation) WithTx(tx *gorm.DB) *MGetKnowledgeBaseByIDsOperation {
	op.tx = tx
	return op
}

func (op *MGetKnowledgeBaseByIDsOperation) Exec() ([]*model.KnowledgeBase, error) {
	if len(op.ids) == 0 {
		return []*model.KnowledgeBase{}, nil
	}
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	var kbs []*model.KnowledgeBase
	if err := db.Model(&model.KnowledgeBase{}).Where("id IN ?", op.ids).Find(&kbs).Error; err != nil {
		return nil, err
	}
	return kbs, nil
}
