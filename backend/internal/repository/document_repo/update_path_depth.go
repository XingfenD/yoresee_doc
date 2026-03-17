package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentUpdatePathDepthOperation struct {
	repo     *DocumentRepository
	docID    int64
	parentID int64
	tx       *gorm.DB
}

func (r *DocumentRepository) UpdatePathDepth(docID, parentID int64) *DocumentUpdatePathDepthOperation {
	return &DocumentUpdatePathDepthOperation{
		repo:     r,
		docID:    docID,
		parentID: parentID,
	}
}

func (op *DocumentUpdatePathDepthOperation) WithTx(tx *gorm.DB) *DocumentUpdatePathDepthOperation {
	op.tx = tx
	return op
}

func (op *DocumentUpdatePathDepthOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	if op.parentID == 0 {
		return op.tx.Exec(`
			UPDATE document_metas
			SET path = (id::text)::ltree, depth = 0
			WHERE id = ? AND deleted_at IS NULL
		`, op.docID).Error
	}

	type pathDepth struct {
		Path  string
		Depth int
	}
	var parent pathDepth
	if err := op.tx.Model(&model.Document{}).
		Select("path, depth").
		Where("id = ? AND deleted_at IS NULL", op.parentID).
		Take(&parent).Error; err != nil {
		return err
	}

	return op.tx.Exec(`
		UPDATE document_metas
		SET path = (?::ltree) || (id::text)::ltree, depth = ?
		WHERE id = ? AND deleted_at IS NULL
	`, parent.Path, parent.Depth+1, op.docID).Error
}
