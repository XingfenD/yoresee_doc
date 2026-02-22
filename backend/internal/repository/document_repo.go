package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

// DocumentRepository 文档仓库
type DocumentRepository struct{}

var DocumentRepo = &DocumentRepository{}

// DocumentGetByExternalIDOperation 根据ExternalID获取文档操作
type DocumentGetByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	tx         *gorm.DB
}

func (r *DocumentRepository) GetByExternalID(externalID string) *DocumentGetByExternalIDOperation {
	return &DocumentGetByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *DocumentGetByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentGetByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetByExternalIDOperation) Exec() (*model.DocumentMeta, error) {
	var document model.DocumentMeta
	var err error

	if op.tx != nil {
		err = op.tx.First(&document, "external_id = ?", op.externalID).Error
	} else {
		err = storage.DB.First(&document, "external_id = ?", op.externalID).Error
	}

	return &document, err
}

// DocumentGetContentOperation 获取文档内容操作
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
	var version model.DocumentVersion
	var content model.Content
	var err error

	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	// 获取最新版本
	err = db.Where("document_id = ?", op.documentID).Order("version DESC").First(&version).Error
	if err != nil {
		return "", err
	}

	// 获取内容
	err = db.First(&content, version.ContentID).Error
	if err != nil {
		return "", err
	}

	return content.Content, nil
}
