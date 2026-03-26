package comment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CommentUpdateContentOperation struct {
	repo    *CommentRepository
	id      int64
	content string
	tx      *gorm.DB
}

func (r *CommentRepository) UpdateContentByID(id int64, content string) *CommentUpdateContentOperation {
	return &CommentUpdateContentOperation{
		repo:    r,
		id:      id,
		content: content,
	}
}

func (op *CommentUpdateContentOperation) WithTx(tx *gorm.DB) *CommentUpdateContentOperation {
	op.tx = tx
	return op
}

func (op *CommentUpdateContentOperation) Exec() error {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Model(&model.DocumentComment{}).Where("id = ?", op.id).Update("content", op.content).Error
}
