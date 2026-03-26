package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UpsertRecentDocumentOperation struct {
	repo *DocumentRepository
	m    *model.RecentDocument
	tx   *gorm.DB
}

func (r *DocumentRepository) UpsertRecentDocument(m *model.RecentDocument) *UpsertRecentDocumentOperation {
	return &UpsertRecentDocumentOperation{
		repo: r,
		m:    m,
	}
}

func (op *UpsertRecentDocumentOperation) WithTx(tx *gorm.DB) *UpsertRecentDocumentOperation {
	op.tx = tx
	return op
}

func (op *UpsertRecentDocumentOperation) Exec() error {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "document_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"accessed_at": op.m.AccessedAt}),
	}).Create(op.m).Error
}
