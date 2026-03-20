package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type TemplateMGetByIDsOperation struct {
	repo *TemplateRepository
	ids  []int64
	tx   *gorm.DB
}

func (r *TemplateRepository) MGetByIDs(ids []int64) *TemplateMGetByIDsOperation {
	return &TemplateMGetByIDsOperation{
		repo: r,
		ids:  ids,
	}
}

func (op *TemplateMGetByIDsOperation) WithTx(tx *gorm.DB) *TemplateMGetByIDsOperation {
	op.tx = tx
	return op
}

func (op *TemplateMGetByIDsOperation) Exec() ([]*model.Template, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	var templates []*model.Template
	if len(op.ids) == 0 {
		return templates, nil
	}
	if err := db.Where("id IN ?", op.ids).Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
