package comment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CommentGetByExternalIDOperation struct {
	repo       *CommentRepository
	externalID string
	tx         *gorm.DB
}

func (r *CommentRepository) GetByExternalID(externalID string) *CommentGetByExternalIDOperation {
	return &CommentGetByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *CommentGetByExternalIDOperation) WithTx(tx *gorm.DB) *CommentGetByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *CommentGetByExternalIDOperation) Exec() (*model.DocumentComment, error) {
	var item model.DocumentComment
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	if err := db.First(&item, "external_id = ?", op.externalID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
