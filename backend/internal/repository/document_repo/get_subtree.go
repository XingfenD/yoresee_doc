package document_repo

import (
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentGetSubtreeOperation struct {
	repo         *DocumentRepository
	rootParentID int64
	knowledgeID  *int64
	depth        *int
	tx           *gorm.DB
}

func (r *DocumentRepository) GetSubtree(rootParentID int64) *DocumentGetSubtreeOperation {
	return &DocumentGetSubtreeOperation{
		repo:         r,
		rootParentID: rootParentID,
	}
}

func (op *DocumentGetSubtreeOperation) WithTx(tx *gorm.DB) *DocumentGetSubtreeOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetSubtreeOperation) WithKnowledgeID(knowledgeID *int64) *DocumentGetSubtreeOperation {
	op.knowledgeID = knowledgeID
	return op
}

func (op *DocumentGetSubtreeOperation) WithDepth(depth int) *DocumentGetSubtreeOperation {
	op.depth = &depth
	return op
}

func (op *DocumentGetSubtreeOperation) Exec() ([]model.Document, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	var documents []model.Document

	if op.depth != nil && *op.depth == 0 {
		return documents, nil
	}

	type pathDepth struct {
		Path  string
		Depth int
	}
	var root pathDepth
	err := db.Model(&model.Document{}).
		Select("path, depth").
		Where("id = ?", op.rootParentID).
		Take(&root).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return documents, nil
		}
		return nil, err
	}

	query := `
		SELECT *
		FROM document_metas
		WHERE deleted_at IS NULL
			AND id <> ?
			AND path <@ ?
	`
	args := []interface{}{op.rootParentID, root.Path}
	if op.knowledgeID != nil {
		query += " AND knowledge_id = ?"
		args = append(args, *op.knowledgeID)
	}
	if op.depth != nil {
		maxDepth := root.Depth + *op.depth
		query += " AND depth <= ?"
		args = append(args, maxDepth)
	}
	query += " ORDER BY depth, created_at"

	err = db.Raw(query, args...).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}
