package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentVersionRepository struct{}

var DocumentVersionRepo = &DocumentVersionRepository{}

type DocumentVersionCreateOperation struct {
	repo       *DocumentVersionRepository
	docVersion *model.DocumentVersion
	tx         *gorm.DB
}

func (repo *DocumentVersionRepository) Create(docVersion *model.DocumentVersion) *DocumentVersionCreateOperation {
	return &DocumentVersionCreateOperation{
		repo:       repo,
		docVersion: docVersion,
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
		Where("document_id = ?", op.docVersion.DocumentID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return err
	}

	op.docVersion.Version = maxVersion + 1

	return op.tx.Create(op.docVersion).Error
}

type DocumentVersionGetLastedOperation struct {
	repo  *DocumentVersionRepository
	docID int64
	tx    *gorm.DB
}

func (r *DocumentVersionRepository) GetLasted(docID int64) *DocumentVersionGetLastedOperation {
	return &DocumentVersionGetLastedOperation{
		repo:  r,
		docID: docID,
	}
}

func (op *DocumentVersionGetLastedOperation) WithTx(tx *gorm.DB) *DocumentVersionGetLastedOperation {
	op.tx = tx
	return op
}

func (op *DocumentVersionGetLastedOperation) Exec() (int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	var maxVersion int64
	err := op.tx.Model(&model.DocumentVersion{}).
		Where("document_id = ?", op.docID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}

	return maxVersion, nil
}
