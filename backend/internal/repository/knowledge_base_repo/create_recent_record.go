package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CreateRecentKnowledgeBaseOperation struct {
	repo *KnowledgeBaseRepository
	m    *model.RecentKnowledgeBase
	tx   *gorm.DB
}

func (r *KnowledgeBaseRepository) CreateRecentKnowledgeBase(m *model.RecentKnowledgeBase) *CreateRecentKnowledgeBaseOperation {
	return &CreateRecentKnowledgeBaseOperation{
		repo: r,
		m:    m,
	}
}

func (op *CreateRecentKnowledgeBaseOperation) WithTx(tx *gorm.DB) {
	op.tx = tx
}

func (op *CreateRecentKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "knowledge_base_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"accessed_at"}),
	}).Create(op.m).Error
}
