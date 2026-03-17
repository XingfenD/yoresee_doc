package document_version_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
