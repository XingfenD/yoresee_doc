package comment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type CommentCreateOperation struct {
	repo *CommentRepository
	item *model.DocumentComment
	tx   *gorm.DB
}

func (r *CommentRepository) Create(item *model.DocumentComment) *CommentCreateOperation {
	return &CommentCreateOperation{
		repo: r,
		item: item,
	}
}

func (op *CommentCreateOperation) WithTx(tx *gorm.DB) *CommentCreateOperation {
	op.tx = tx
	return op
}

func (op *CommentCreateOperation) Exec() error {
	if op.item == nil {
		return nil
	}
	db := op.repo.db
	if op.tx != nil {
		db = op.tx
	}
	return db.Create(op.item).Error
}
