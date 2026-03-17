package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentMoveSubtreeOperation struct {
	repo        *DocumentRepository
	docID       int64
	newParentID int64
	tx          *gorm.DB
}

func (r *DocumentRepository) MoveSubtree(docID, newParentID int64) *DocumentMoveSubtreeOperation {
	return &DocumentMoveSubtreeOperation{
		repo:        r,
		docID:       docID,
		newParentID: newParentID,
	}
}

func (op *DocumentMoveSubtreeOperation) WithTx(tx *gorm.DB) *DocumentMoveSubtreeOperation {
	op.tx = tx
	return op
}

func (op *DocumentMoveSubtreeOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	type pathDepth struct {
		Path  string
		Depth int
	}
	var old pathDepth
	if err := op.tx.Model(&model.Document{}).
		Select("path, depth").
		Where("id = ? AND deleted_at IS NULL", op.docID).
		Take(&old).Error; err != nil {
		return err
	}

	var newPath string
	var newDepth int
	if op.newParentID == 0 {
		newDepth = 0
	} else {
		var parent pathDepth
		if err := op.tx.Model(&model.Document{}).
			Select("path, depth").
			Where("id = ? AND deleted_at IS NULL", op.newParentID).
			Take(&parent).Error; err != nil {
			return err
		}
		newPath = parent.Path
		newDepth = parent.Depth + 1
	}

	depthDelta := newDepth - old.Depth

	if op.newParentID == 0 {
		return op.tx.Exec(`
			UPDATE document_metas
			SET path = (?::text)::ltree || COALESCE(subpath(path, nlevel(?::ltree)), ''::ltree),
				depth = depth + ?
			WHERE path <@ ?::ltree AND deleted_at IS NULL
		`, op.docID, old.Path, depthDelta, old.Path).Error
	}

	return op.tx.Exec(`
		UPDATE document_metas
		SET path = (?::ltree) || (?::text)::ltree || COALESCE(subpath(path, nlevel(?::ltree)), ''::ltree),
			depth = depth + ?
		WHERE path <@ ?::ltree AND deleted_at IS NULL
	`, newPath, op.docID, old.Path, depthDelta, old.Path).Error
}
