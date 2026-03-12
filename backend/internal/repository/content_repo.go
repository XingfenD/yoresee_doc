package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ContentRepository struct{}

var ContentRepo = &ContentRepository{}

type ContentCreateOperation struct {
	repo    *ContentRepository
	content string
	tx      *gorm.DB
}

func (r *ContentRepository) Create(content string) *ContentCreateOperation {
	return &ContentCreateOperation{
		repo:    r,
		content: content,
	}
}

func (op *ContentCreateOperation) WithTx(tx *gorm.DB) *ContentCreateOperation {
	op.tx = tx
	return op
}

func (op *ContentCreateOperation) Exec() (int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	m := &model.Content{
		Content: op.content,
	}
	if err := op.tx.Create(m).Error; err != nil {
		return 0, err
	}
	return m.ID, nil
}

type ContentUpdateOperation struct {
	repo         *ContentRepository
	contentModel *model.Content
	tx           *gorm.DB
}

func (r *ContentRepository) Update(contentModel *model.Content) *ContentUpdateOperation {
	return &ContentUpdateOperation{
		repo:         r,
		contentModel: contentModel,
	}
}

func (op *ContentUpdateOperation) WithTx(tx *gorm.DB) *ContentUpdateOperation {
	op.tx = tx
	return op
}

func (op *ContentUpdateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Save(op.contentModel).Error
}
