package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentGetSubtreeByKnowledgeIDOperation struct {
	repo        *DocumentRepository
	knowledgeID int64
	depth       *int
	tx          *gorm.DB
}

func (r *DocumentRepository) GetSubtreeByKnowledgeID(knowledgeID int64) *DocumentGetSubtreeByKnowledgeIDOperation {
	return &DocumentGetSubtreeByKnowledgeIDOperation{
		repo:        r,
		knowledgeID: knowledgeID,
	}
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) WithTx(tx *gorm.DB) *DocumentGetSubtreeByKnowledgeIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) WithDepth(depth int) *DocumentGetSubtreeByKnowledgeIDOperation {
	op.depth = &depth
	return op
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) Exec() ([]model.Document, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	var documents []model.Document

	if op.depth != nil && *op.depth == 0 {
		return documents, nil
	}

	query := `
		SELECT *
		FROM document_metas
		WHERE deleted_at IS NULL
			AND knowledge_id = ?
	`

	args := []interface{}{op.knowledgeID}
	if op.depth != nil {
		query += " AND depth <= ?"
		args = append(args, *op.depth)
	}
	query += " ORDER BY depth, created_at"

	err := db.Raw(query, args...).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}
