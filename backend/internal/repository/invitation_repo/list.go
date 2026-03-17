package invitation_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
