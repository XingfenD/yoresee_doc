package comment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CommentDeleteOperation struct {
	repo *CommentRepository
	id   int64
	tx   *gorm.DB
}

func (r *CommentRepository) Delete(id int64) *CommentDeleteOperation {
	return &CommentDeleteOperation{
		repo: r,
		id:   id,
	}
}

func (op *CommentDeleteOperation) WithTx(tx *gorm.DB) *CommentDeleteOperation {
	op.tx = tx
	return op
}

func (op *CommentDeleteOperation) Exec() error {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Delete(&model.DocumentComment{}, op.id).Error
}
