package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type KnowledgeBaseRepository struct{}

var KnowledgeBaseRepo = &KnowledgeBaseRepository{}

type ListKnowledgeBaseOperation struct {
	repo  *KnowledgeBaseRepository
	model *model.KnowledgeBase
	tx    *gorm.DB
}

func (r *KnowledgeBaseRepository) List(m *model.KnowledgeBase) (op *ListKnowledgeBaseOperation) {
	return &ListKnowledgeBaseOperation{
		repo:  KnowledgeBaseRepo,
		model: m,
	}
}

func (op *ListKnowledgeBaseOperation) WithTx(tx *gorm.DB) *ListKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *ListKnowledgeBaseOperation) Exec() (kbs []*model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err = op.tx.Model(op.model).Find(&kbs).Error
	return
}

type KnowledgeBaseGetIDByExternalIDOperation struct {
	repo       *KnowledgeBaseRepository
	externalID string
	tx         *gorm.DB
}

func (r *KnowledgeBaseRepository) GetIDByExternalID(externalID string) (op *KnowledgeBaseGetIDByExternalIDOperation) {
	return &KnowledgeBaseGetIDByExternalIDOperation{
		repo:       KnowledgeBaseRepo,
		externalID: externalID,
	}
}

func (op *KnowledgeBaseGetIDByExternalIDOperation) WithTx(tx *gorm.DB) *KnowledgeBaseGetIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *KnowledgeBaseGetIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.First(&model.KnowledgeBase{}, "external_id = ?", op.externalID).Pluck("id", &id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

type KnowledgeBaseGetByIDOperation struct {
	repo *KnowledgeBaseRepository
	id   int64
	tx   *gorm.DB
}

func (r *KnowledgeBaseRepository) GetByID(id int64) (op *KnowledgeBaseGetByIDOperation) {
	return &KnowledgeBaseGetByIDOperation{
		repo: KnowledgeBaseRepo,
		id:   id,
	}
}

func (op *KnowledgeBaseGetByIDOperation) WithTx(tx *gorm.DB) *KnowledgeBaseGetByIDOperation {
	op.tx = tx
	return op
}

func (op *KnowledgeBaseGetByIDOperation) Exec() (knowledgeBase *model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err = op.tx.First(knowledgeBase, "id = ?", op.id).Error
	return
}

type GetKnowledgeBaseByExternalIDOperation struct {
	repo       *KnowledgeBaseRepository
	externalID string
	tx         *gorm.DB
}

func (r *KnowledgeBaseRepository) GetByExternalID(externalID string) (op *GetKnowledgeBaseByExternalIDOperation) {
	return &GetKnowledgeBaseByExternalIDOperation{
		repo:       KnowledgeBaseRepo,
		externalID: externalID,
	}
}

func (op *GetKnowledgeBaseByExternalIDOperation) WithTx(tx *gorm.DB) *GetKnowledgeBaseByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *GetKnowledgeBaseByExternalIDOperation) Exec() (knowledgeBase *model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err = op.tx.First(knowledgeBase, "external_id = ?", op.externalID).Error
	return
}

type CreateKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Create(knowledgeBase *model.KnowledgeBase) (op *CreateKnowledgeBaseOperation) {
	return &CreateKnowledgeBaseOperation{
		repo:          KnowledgeBaseRepo,
		knowledgeBase: knowledgeBase,
	}
}

func (op *CreateKnowledgeBaseOperation) WithTx(tx *gorm.DB) *CreateKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *CreateKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.Create(op.knowledgeBase).Error
	return err
}

type DeleteKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	knowledgeBase *model.KnowledgeBase
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) Delete(knowledgeBase *model.KnowledgeBase) (op *DeleteKnowledgeBaseOperation) {
	return &DeleteKnowledgeBaseOperation{
		repo:          KnowledgeBaseRepo,
		knowledgeBase: knowledgeBase,
	}
}

func (op *DeleteKnowledgeBaseOperation) WithTx(tx *gorm.DB) *DeleteKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *DeleteKnowledgeBaseOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.Delete(op.knowledgeBase).Error
	return err
}
