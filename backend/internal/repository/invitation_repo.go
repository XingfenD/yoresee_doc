package repository

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationRepository struct{}

var InvitationRepo = &InvitationRepository{}

type InvitationCreateOperation struct {
	repo       *InvitationRepository
	invitation *model.Invitation
	tx         *gorm.DB
}

func (r *InvitationRepository) Create(invitation *model.Invitation) *InvitationCreateOperation {
	return &InvitationCreateOperation{
		repo:       r,
		invitation: invitation,
	}
}

func (op *InvitationCreateOperation) WithTx(tx *gorm.DB) *InvitationCreateOperation {
	op.tx = tx
	return op
}

func (op *InvitationCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.invitation).Error
	}
	return storage.DB.Create(op.invitation).Error
}

type InvitationGetByCodeOperation struct {
	repo *InvitationRepository
	code string
	tx   *gorm.DB
}

func (r *InvitationRepository) GetByCode(code string) *InvitationGetByCodeOperation {
	return &InvitationGetByCodeOperation{
		repo: r,
		code: code,
	}
}

func (op *InvitationGetByCodeOperation) WithTx(tx *gorm.DB) *InvitationGetByCodeOperation {
	op.tx = tx
	return op
}

func (op *InvitationGetByCodeOperation) Exec() (*model.Invitation, error) {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Where("code = ?", op.code).First(&invitation).Error
	} else {
		err = storage.DB.Where("code = ?", op.code).First(&invitation).Error
	}

	return &invitation, err
}

type InvitationGetByIDOperation struct {
	repo *InvitationRepository
	id   int64
	tx   *gorm.DB
}

func (r *InvitationRepository) GetByID(id int64) *InvitationGetByIDOperation {
	return &InvitationGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *InvitationGetByIDOperation) WithTx(tx *gorm.DB) *InvitationGetByIDOperation {
	op.tx = tx
	return op
}

func (op *InvitationGetByIDOperation) Exec() (*model.Invitation, error) {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.First(&invitation, op.id).Error
	} else {
		err = storage.DB.First(&invitation, op.id).Error
	}

	return &invitation, err
}

type InvitationListOperation struct {
	repo           *InvitationRepository
	m              *model.Invitation
	creatorID      *int64
	maxUsedCnt     *int64
	expiresAtStart *string
	expiresAtEnd   *string
	createdAtStart *string
	createdAtEnd   *string
	disabled       *bool
	sortField      string
	sortDesc       bool
	page           int
	pageSize       int
	tx             *gorm.DB
}

func (r *InvitationRepository) List(m *model.Invitation) *InvitationListOperation {
	return &InvitationListOperation{
		repo: r,
		m:    m,
	}
}

func (op *InvitationListOperation) WithTx(tx *gorm.DB) *InvitationListOperation {
	op.tx = tx
	return op
}

func (op *InvitationListOperation) WithCreatorID(creatorID *int64) *InvitationListOperation {
	op.creatorID = creatorID
	return op
}

func (op *InvitationListOperation) WithMaxUsedCnt(maxUsedCnt *int64) *InvitationListOperation {
	op.maxUsedCnt = maxUsedCnt
	return op
}

func (op *InvitationListOperation) WithExpiresAtRange(start, end *string) *InvitationListOperation {
	op.expiresAtStart = start
	op.expiresAtEnd = end
	return op
}

func (op *InvitationListOperation) WithCreatedAtRange(start, end *string) *InvitationListOperation {
	op.createdAtStart = start
	op.createdAtEnd = end
	return op
}

func (op *InvitationListOperation) WithDisabled(disabled *bool) *InvitationListOperation {
	op.disabled = disabled
	return op
}

func (op *InvitationListOperation) WithSort(field string, desc bool) *InvitationListOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *InvitationListOperation) WithPagination(page, pageSize int) *InvitationListOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *InvitationListOperation) Exec() ([]model.Invitation, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	dbQuery := db.Model(&model.Invitation{})

	if op.m != nil {
		dbQuery = dbQuery.Where(op.m)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("created_by = ?", *op.creatorID)
	}

	if op.maxUsedCnt != nil {
		dbQuery = dbQuery.Where("max_used_cnt = ?", *op.maxUsedCnt)
	}

	if op.expiresAtStart != nil && *op.expiresAtStart != "" {
		dbQuery = dbQuery.Where("expires_at >= ?", *op.expiresAtStart)
	}
	if op.expiresAtEnd != nil && *op.expiresAtEnd != "" {
		dbQuery = dbQuery.Where("expires_at <= ?", *op.expiresAtEnd)
	}
	if op.createdAtStart != nil && *op.createdAtStart != "" {
		dbQuery = dbQuery.Where("created_at >= ?", *op.createdAtStart)
	}
	if op.createdAtEnd != nil && *op.createdAtEnd != "" {
		dbQuery = dbQuery.Where("created_at <= ?", *op.createdAtEnd)
	}
	if op.disabled != nil {
		dbQuery = dbQuery.Where("disabled = ?", *op.disabled)
	}

	sortField := op.sortField
	if sortField == "" {
		sortField = "created_at"
	}
	orderDirection := "ASC"
	if op.sortDesc {
		orderDirection = "DESC"
	}
	dbQuery = dbQuery.Order(sortField + " " + orderDirection)

	page := op.page
	pageSize := op.pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Offset(offset).Limit(pageSize)

	var invitations []model.Invitation
	err := dbQuery.Find(&invitations).Error
	if err != nil {
		return nil, err
	}

	return invitations, nil
}

func (op *InvitationListOperation) ExecWithTotal() ([]model.Invitation, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	dbQuery := db.Model(&model.Invitation{})

	if op.m != nil {
		dbQuery = dbQuery.Where(op.m)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("created_by = ?", *op.creatorID)
	}

	if op.maxUsedCnt != nil {
		dbQuery = dbQuery.Where("max_used_cnt = ?", *op.maxUsedCnt)
	}

	if op.expiresAtStart != nil && *op.expiresAtStart != "" {
		dbQuery = dbQuery.Where("expires_at >= ?", *op.expiresAtStart)
	}
	if op.expiresAtEnd != nil && *op.expiresAtEnd != "" {
		dbQuery = dbQuery.Where("expires_at <= ?", *op.expiresAtEnd)
	}
	if op.createdAtStart != nil && *op.createdAtStart != "" {
		dbQuery = dbQuery.Where("created_at >= ?", *op.createdAtStart)
	}
	if op.createdAtEnd != nil && *op.createdAtEnd != "" {
		dbQuery = dbQuery.Where("created_at <= ?", *op.createdAtEnd)
	}
	if op.disabled != nil {
		dbQuery = dbQuery.Where("disabled = ?", *op.disabled)
	}

	var total int64
	err := dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	sortField := op.sortField
	if sortField == "" {
		sortField = "created_at"
	}
	orderDirection := "ASC"
	if op.sortDesc {
		orderDirection = "DESC"
	}
	dbQuery = dbQuery.Order(sortField + " " + orderDirection)

	page := op.page
	pageSize := op.pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Offset(offset).Limit(pageSize)

	var invitations []model.Invitation
	err = dbQuery.Find(&invitations).Error
	if err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}

// type InvitationListByCreatedByOperation struct {
// 	repo      *InvitationRepository
// 	createdBy int64
// 	tx        *gorm.DB
// }

// func (r *InvitationRepository) ListByCreatedBy(createdBy int64) *InvitationListByCreatedByOperation {
// 	return &InvitationListByCreatedByOperation{
// 		repo:      r,
// 		createdBy: createdBy,
// 	}
// }

// func (op *InvitationListByCreatedByOperation) WithTx(tx *gorm.DB) *InvitationListByCreatedByOperation {
// 	op.tx = tx
// 	return op
// }

// func (op *InvitationListByCreatedByOperation) Exec() ([]model.Invitation, error) {
// 	var invitations []model.Invitation
// 	var err error

// 	if op.tx != nil {
// 		err = op.tx.Where("created_by = ?", op.createdBy).Find(&invitations).Error
// 	} else {
// 		err = storage.DB.Where("created_by = ?", op.createdBy).Find(&invitations).Error
// 	}

// 	return invitations, err
// }

type InvitationDeleteOperation struct {
	repo *InvitationRepository
	id   int64
	tx   *gorm.DB
}

func (r *InvitationRepository) Delete(id int64) *InvitationDeleteOperation {
	return &InvitationDeleteOperation{
		repo: r,
		id:   id,
	}
}

func (op *InvitationDeleteOperation) WithTx(tx *gorm.DB) *InvitationDeleteOperation {
	op.tx = tx
	return op
}

func (op *InvitationDeleteOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	return op.tx.Delete(&model.Invitation{}, op.id).Error
}

type InvitationUpdateOperation struct {
	repo       *InvitationRepository
	invitation *model.Invitation
	tx         *gorm.DB
}

func (r *InvitationRepository) Update(invitation *model.Invitation) *InvitationUpdateOperation {
	return &InvitationUpdateOperation{
		repo:       r,
		invitation: invitation,
	}
}

func (op *InvitationUpdateOperation) WithTx(tx *gorm.DB) *InvitationUpdateOperation {
	op.tx = tx
	return op
}

func (op *InvitationUpdateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Save(op.invitation).Error
	}
	return storage.DB.Save(op.invitation).Error
}

type InvitationValidateAndUseOperation struct {
	repo *InvitationRepository
	code string
	tx   *gorm.DB
}

func (r *InvitationRepository) ValidateAndUse(code string) *InvitationValidateAndUseOperation {
	return &InvitationValidateAndUseOperation{
		repo: r,
		code: code,
	}
}

func (op *InvitationValidateAndUseOperation) WithTx(tx *gorm.DB) *InvitationValidateAndUseOperation {
	op.tx = tx
	return op
}

func (op *InvitationValidateAndUseOperation) Exec() error {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Where("code = ? AND is_used = ?", op.code, false).First(&invitation).Error
	} else {
		err = storage.DB.Where("code = ? AND is_used = ?", op.code, false).First(&invitation).Error
	}

	if err != nil {
		return err
	}

	if op.tx != nil {
		return op.tx.Model(&invitation).Updates(map[string]interface{}{
			"is_used": true,
			"used_at": time.Now(),
		}).Error
	}

	return storage.DB.Model(&invitation).Updates(map[string]interface{}{
		"is_used": true,
		"used_at": time.Now(),
	}).Error
}
