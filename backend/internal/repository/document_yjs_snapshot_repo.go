package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocumentYjsSnapshotRepository struct{}

var DocumentYjsSnapshotRepo = &DocumentYjsSnapshotRepository{}

type DocumentYjsSnapshotGetOperation struct {
	repo  *DocumentYjsSnapshotRepository
	docID int64
	tx    *gorm.DB
}

func (r *DocumentYjsSnapshotRepository) GetByDocID(docID int64) *DocumentYjsSnapshotGetOperation {
	return &DocumentYjsSnapshotGetOperation{
		repo:  r,
		docID: docID,
	}
}

func (op *DocumentYjsSnapshotGetOperation) WithTx(tx *gorm.DB) *DocumentYjsSnapshotGetOperation {
	op.tx = tx
	return op
}

func (op *DocumentYjsSnapshotGetOperation) Exec() (*model.DocumentYjsSnapshot, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	var snapshot model.DocumentYjsSnapshot
	if err := op.tx.First(&snapshot, "doc_id = ?", op.docID).Error; err != nil {
		return nil, err
	}
	return &snapshot, nil
}

type DocumentYjsSnapshotSaveOperation struct {
	repo  *DocumentYjsSnapshotRepository
	docID int64
	state []byte
	tx    *gorm.DB
}

func (r *DocumentYjsSnapshotRepository) Save(docID int64, state []byte) *DocumentYjsSnapshotSaveOperation {
	return &DocumentYjsSnapshotSaveOperation{
		repo:  r,
		docID: docID,
		state: state,
	}
}

func (op *DocumentYjsSnapshotSaveOperation) WithTx(tx *gorm.DB) *DocumentYjsSnapshotSaveOperation {
	op.tx = tx
	return op
}

func (op *DocumentYjsSnapshotSaveOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	snapshot := &model.DocumentYjsSnapshot{
		DocID:    op.docID,
		YjsState: op.state,
		Version:  1,
	}

	return op.tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "doc_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"yjs_state":  op.state,
			"version":    gorm.Expr("documents_yjs_snapshot.version + 1"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(snapshot).Error
}
