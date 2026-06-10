package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type TemplateGetByIDOperation struct {
	repo *TemplateRepository
	id   int64
	tx   *gorm.DB
}

func (r *TemplateRepository) GetByID(id int64) (op *TemplateGetByIDOperation) {
	return &TemplateGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *TemplateGetByIDOperation) WithTx(tx *gorm.DB) *TemplateGetByIDOperation {
	op.tx = tx
	return op
}

func (op *TemplateGetByIDOperation) Exec() (template *model.Template, err error) {
	if op.tx == nil {
		op.tx = op.repo.db
	}
	template = &model.Template{}
	err = op.tx.First(template, "id = ?", op.id).Error
	return
}
