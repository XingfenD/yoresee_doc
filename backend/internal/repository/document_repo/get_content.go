package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type DocumentGetContentOperation struct {
	repo       *DocumentRepository
	documentID int64
	tx         *gorm.DB
}

func (r *DocumentRepository) GetContent(documentID int64) *DocumentGetContentOperation {
	return &DocumentGetContentOperation{
		repo:       r,
		documentID: documentID,
	}
}

func (op *DocumentGetContentOperation) WithTx(tx *gorm.DB) *DocumentGetContentOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetContentOperation) Exec() (string, error) {
	var docMeta model.Document
	var err error

	db := op.repo.db
	if op.tx != nil {
		db = op.tx
	}

	err = db.Where("id = ?", op.documentID).First(&docMeta).Error
	if err != nil {
		return "", err
	}

	return docMeta.Content, nil
}
