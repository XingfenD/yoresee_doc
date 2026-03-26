package comment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CommentListByDocumentOperation struct {
	repo       *CommentRepository
	documentID int64
	page       int
	pageSize   int
	tx         *gorm.DB
}

func (r *CommentRepository) ListByDocument(documentID int64) *CommentListByDocumentOperation {
	return &CommentListByDocumentOperation{
		repo:       r,
		documentID: documentID,
	}
}

func (op *CommentListByDocumentOperation) WithTx(tx *gorm.DB) *CommentListByDocumentOperation {
	op.tx = tx
	return op
}

func (op *CommentListByDocumentOperation) WithPagination(page, pageSize int) *CommentListByDocumentOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *CommentListByDocumentOperation) ExecWithTotal() ([]model.DocumentComment, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	query := db.Model(&model.DocumentComment{}).
		Where("document_id = ?", op.documentID).
		Where("anchor_id IS NOT NULL AND anchor_id <> ''")
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := op.page
	pageSize := op.pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var items []model.DocumentComment
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
