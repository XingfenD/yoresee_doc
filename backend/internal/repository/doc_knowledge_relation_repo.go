package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
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

type DocKBRelationCreateOperation struct {
	repo    *DocKnowledgeRelationRepository
	docID   int64
	kbID    *int64
	ownerID *int64
	tx      *gorm.DB
}

func (r *DocKnowledgeRelationRepository) CreateDocKBRelation(docID int64) *DocKBRelationCreateOperation {
	return &DocKBRelationCreateOperation{
		repo:  r,
		docID: docID,
	}
}

func (op *DocKBRelationCreateOperation) WithKnowledgeID(kbID *int64) *DocKBRelationCreateOperation {
	op.kbID = kbID
	return op
}

func (op *DocKBRelationCreateOperation) WithOwnerID(ownerID *int64) *DocKBRelationCreateOperation {
	op.ownerID = ownerID
	return op
}

func (op *DocKBRelationCreateOperation) WithTx(tx *gorm.DB) *DocKBRelationCreateOperation {
	op.tx = tx
	return op
}

func (op *DocKBRelationCreateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	if op.kbID != nil && op.ownerID != nil {
		return status.StatusInternalParamsError
	}

	relation := model.DocKnowledgeRelation{
		DocumentID:  op.docID,
		KnowledgeID: op.kbID,
		OwnerID:     op.ownerID,
	}

	return op.tx.FirstOrCreate(&relation, model.DocKnowledgeRelation{
		DocumentID: relation.DocumentID,
	}).Error
}
