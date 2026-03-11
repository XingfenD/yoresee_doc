package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentVersionRepository struct{}

var DocumentVersionRepo = &DocumentVersionRepository{}

type DocumentVersionCreateOperation struct {
	repo  *DocumentVersionRepository
	docID int64
	tx    *gorm.DB
}

func (repo *DocumentVersionRepository) Create(docID int64) *DocumentVersionCreateOperation {
	return &DocumentVersionCreateOperation{
		repo:  repo,
		docID: docID,
	}
}

func (op *DocumentVersionCreateOperation) WithTx(tx *gorm.DB) *DocumentVersionCreateOperation {
	op.tx = tx
	return op
}

func (op *DocumentVersionCreateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	var maxVersion int
	err := op.tx.Model(&model.DocumentVersion{}).
		Where("document_id = ?", op.docID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return err
	}

	version := maxVersion + 1

	return op.tx.Create(&model.DocumentVersion{
		DocumentID: op.docID,
		Version:    version,
	}).Error
}
