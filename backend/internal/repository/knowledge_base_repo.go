package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type KnowledgeBaseRepository struct{}

var KnowledgeBaseRepo = &KnowledgeBaseRepository{}

type ListKnowledgeBaseOperation struct {
	repo      *KnowledgeBaseRepository
	model     *model.KnowledgeBase
	creatorID *int64
	isPublic  *bool
	sortField string
	sortDesc  bool
	page      int
	pageSize  int
	tx        *gorm.DB
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

func (op *ListKnowledgeBaseOperation) WithCreatorID(creatorID *int64) *ListKnowledgeBaseOperation {
	op.creatorID = creatorID
	return op
}

func (op *ListKnowledgeBaseOperation) WithIsPublic(isPublic *bool) *ListKnowledgeBaseOperation {
	op.isPublic = isPublic
	return op
}

func (op *ListKnowledgeBaseOperation) WithSort(field string, desc bool) *ListKnowledgeBaseOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *ListKnowledgeBaseOperation) WithPagination(page, pageSize int) *ListKnowledgeBaseOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *ListKnowledgeBaseOperation) Exec() (kbs []*model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	dbQuery := op.tx.Model(op.model)

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("creator_user_id = ?", *op.creatorID)
	}

	if op.isPublic != nil {
		dbQuery = dbQuery.Where("is_public = ?", *op.isPublic)
	}

	orderStr := "created_at DESC"
	if op.sortField != "" {
		if op.sortDesc {
			orderStr = op.sortField + " DESC"
		} else {
			orderStr = op.sortField + " ASC"
		}
	}
	dbQuery = dbQuery.Order(orderStr)

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		dbQuery = dbQuery.Offset(offset).Limit(op.pageSize)
	}

	err = dbQuery.Find(&kbs).Error
	return
}

func (op *ListKnowledgeBaseOperation) ExecWithTotal() (kbs []*model.KnowledgeBase, total int64, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	dbQuery := op.tx.Model(op.model)
	countQuery := op.tx.Model(op.model)

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
		countQuery = countQuery.Where(op.model)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("creator_user_id = ?", *op.creatorID)
		countQuery = countQuery.Where("creator_user_id = ?", *op.creatorID)
	}

	if op.isPublic != nil {
		dbQuery = dbQuery.Where("is_public = ?", *op.isPublic)
		countQuery = countQuery.Where("is_public = ?", *op.isPublic)
	}

	orderStr := "created_at DESC"
	if op.sortField != "" {
		if op.sortDesc {
			orderStr = op.sortField + " DESC"
		} else {
			orderStr = op.sortField + " ASC"
		}
	}
	dbQuery = dbQuery.Order(orderStr)

	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		dbQuery = dbQuery.Offset(offset).Limit(op.pageSize)
	}

	err = dbQuery.Find(&kbs).Error
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
	err := op.tx.Model(&model.KnowledgeBase{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
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

func (op *GetKnowledgeBaseByExternalIDOperation) Exec() (*model.KnowledgeBase, error) {
	var knowledgeBase model.KnowledgeBase
	if op.tx == nil {
		op.tx = storage.DB
	}
	err := op.tx.First(&knowledgeBase, "external_id = ?", op.externalID).Error
	if err != nil {
		return nil, err
	}
	return &knowledgeBase, nil
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

	return op.tx.Create(op.knowledgeBase).Error
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

	return op.tx.Create(op.m).Error
}
