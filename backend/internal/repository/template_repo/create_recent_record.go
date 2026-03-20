package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateRecentTemplateOperation struct {
	repo *TemplateRepository
	m    *model.RecentTemplate
	tx   *gorm.DB
}

func (r *TemplateRepository) CreateRecentTemplate(m *model.RecentTemplate) *CreateRecentTemplateOperation {
	return &CreateRecentTemplateOperation{
		repo: r,
		m:    m,
	}
}

func (op *CreateRecentTemplateOperation) WithTx(tx *gorm.DB) {
	op.tx = tx
}

func (op *CreateRecentTemplateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "template_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"accessed_at"}),
	}).Create(op.m).Error
}
