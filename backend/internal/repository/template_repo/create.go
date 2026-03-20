package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type TemplateCreateOperation struct {
	repo     *TemplateRepository
	template *model.Template
	tx       *gorm.DB
}

func (r *TemplateRepository) Create(template *model.Template) *TemplateCreateOperation {
	return &TemplateCreateOperation{
		repo:     r,
		template: template,
	}
}

func (op *TemplateCreateOperation) WithTx(tx *gorm.DB) *TemplateCreateOperation {
	op.tx = tx
	return op
}

func (op *TemplateCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.template).Error
	}
	return storage.DB.Create(op.template).Error
}
