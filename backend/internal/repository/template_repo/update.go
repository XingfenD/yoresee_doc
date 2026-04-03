package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type TemplateUpdateOperation struct {
	repo         *TemplateRepository
	template     *model.Template
	updateFields map[string]bool
	tx           *gorm.DB
}

func (r *TemplateRepository) Update(template *model.Template) *TemplateUpdateOperation {
	return &TemplateUpdateOperation{
		repo:         r,
		template:     template,
		updateFields: make(map[string]bool),
	}
}

func (op *TemplateUpdateOperation) UpdateName() *TemplateUpdateOperation {
	op.updateFields["name"] = true
	return op
}

func (op *TemplateUpdateOperation) UpdateDescription() *TemplateUpdateOperation {
	op.updateFields["description"] = true
	return op
}

func (op *TemplateUpdateOperation) UpdateScope() *TemplateUpdateOperation {
	op.updateFields["scope"] = true
	return op
}

func (op *TemplateUpdateOperation) UpdateIsPublic() *TemplateUpdateOperation {
	op.updateFields["is_public"] = true
	return op
}

func (op *TemplateUpdateOperation) WithTx(tx *gorm.DB) *TemplateUpdateOperation {
	op.tx = tx
	return op
}

func (op *TemplateUpdateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(op.template)
	if len(op.updateFields) > 0 {
		fields := make([]string, 0, len(op.updateFields))
		for field := range op.updateFields {
			fields = append(fields, field)
		}
		query = query.Select(fields)
	}

	return query.Updates(*op.template).Error
}
