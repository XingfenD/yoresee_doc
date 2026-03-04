package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocKnowledgeRelationRepository struct{}

var DocKnowledgeRelationRepo = &DocKnowledgeRelationRepository{}

type CountDocsByKnowledgeIDOperation struct {
	repo        *DocKnowledgeRelationRepository
	knowledgeID int64
	tx          *gorm.DB
}

func (r *DocKnowledgeRelationRepository) CountDocsByKnowledgeID(knowledgeID int64) *CountDocsByKnowledgeIDOperation {
	return &CountDocsByKnowledgeIDOperation{
		repo:        r,
		knowledgeID: knowledgeID,
	}
}

func (op *CountDocsByKnowledgeIDOperation) WithTx(tx *gorm.DB) *CountDocsByKnowledgeIDOperation {
	op.tx = tx
	return op
}

func (op *CountDocsByKnowledgeIDOperation) Exec() (int64, error) {
	var count int64
	if op.tx == nil {
		op.tx = storage.DB
	}

	err := op.tx.Model(&model.DocKnowledgeRelation{}).
		Where("knowledge_id = ?", op.knowledgeID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
