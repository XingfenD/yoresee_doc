package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type MGetKnowledgeBaseDocumentsCountOperation struct {
	repo             *KnowledgeBaseRepository
	knowledgeBaseIDs []int64
	tx               *gorm.DB
}

func (r *KnowledgeBaseRepository) MGetKnowledgeBaseDocumentsCount(knowledgeBaseIDs []int64) *MGetKnowledgeBaseDocumentsCountOperation {
	return &MGetKnowledgeBaseDocumentsCountOperation{
		repo:             r,
		knowledgeBaseIDs: knowledgeBaseIDs,
	}
}

func (op *MGetKnowledgeBaseDocumentsCountOperation) WithTx(tx *gorm.DB) *MGetKnowledgeBaseDocumentsCountOperation {
	op.tx = tx
	return op
}

func (op *MGetKnowledgeBaseDocumentsCountOperation) Exec() (map[int64]int64, error) {
	result := make(map[int64]int64)

	if len(op.knowledgeBaseIDs) == 0 {
		return result, nil
	}

	if op.tx == nil {
		op.tx = storage.DB
	}

	var counts []struct {
		KnowledgeID int64
		Count       int64
	}

	err := op.tx.Model(&model.Document{}).
		Select("knowledge_id, count(*) as count").
		Where("knowledge_id IN ?", op.knowledgeBaseIDs).
		Group("knowledge_id").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		result[c.KnowledgeID] = c.Count
	}

	for _, id := range op.knowledgeBaseIDs {
		if _, exists := result[id]; !exists {
			result[id] = 0
		}
	}

	return result, nil
}
