package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentCreateOperation struct {
	repo *DocumentRepository
	doc  *model.Document
	tx   *gorm.DB
}

func (r *DocumentRepository) Create(doc *model.Document) *DocumentCreateOperation {
	return &DocumentCreateOperation{
		repo: r,
		doc:  doc,
	}
}

func (op *DocumentCreateOperation) WithTx(tx *gorm.DB) *DocumentCreateOperation {
	op.tx = tx
	return op
}

func (op *DocumentCreateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Create(op.doc).Error
}
