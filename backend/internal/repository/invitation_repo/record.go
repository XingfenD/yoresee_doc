package invitation_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationRecordCreateOperation struct {
	repo   *InvitationRepository
	record *model.InvitationRecord
	tx     *gorm.DB
}

func (r *InvitationRepository) CreateRecord(record *model.InvitationRecord) *InvitationRecordCreateOperation {
	return &InvitationRecordCreateOperation{
		repo:   r,
		record: record,
	}
}

func (op *InvitationRecordCreateOperation) WithTx(tx *gorm.DB) *InvitationRecordCreateOperation {
	op.tx = tx
	return op
}

func (op *InvitationRecordCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.record).Error
	}
	return storage.DB.Create(op.record).Error
}

type InvitationRecordListOperation struct {
	repo        *InvitationRepository
	code        *string
	status      *string
	usedAtStart *string
	usedAtEnd   *string
	creatorID   *int64
	page        int
	pageSize    int
	tx          *gorm.DB
}

func (r *InvitationRepository) ListRecords() *InvitationRecordListOperation {
	return &InvitationRecordListOperation{
		repo: r,
	}
}

func (op *InvitationRecordListOperation) WithTx(tx *gorm.DB) *InvitationRecordListOperation {
	op.tx = tx
	return op
}

func (op *InvitationRecordListOperation) WithCode(code *string) *InvitationRecordListOperation {
	op.code = code
	return op
}

func (op *InvitationRecordListOperation) WithStatus(status *string) *InvitationRecordListOperation {
	op.status = status
	return op
}

func (op *InvitationRecordListOperation) WithUsedAtRange(start, end *string) *InvitationRecordListOperation {
	op.usedAtStart = start
	op.usedAtEnd = end
	return op
}

func (op *InvitationRecordListOperation) WithCreatorID(creatorID *int64) *InvitationRecordListOperation {
	op.creatorID = creatorID
	return op
}

func (op *InvitationRecordListOperation) WithPagination(page, pageSize int) *InvitationRecordListOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *InvitationRecordListOperation) ExecWithTotal() ([]model.InvitationRecord, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	query := db.Model(&model.InvitationRecord{})
	if op.creatorID != nil {
		query = query.Joins("JOIN invitations ON invitations.code = invitation_records.code").
			Where("invitations.created_by = ?", *op.creatorID)
	}
	if op.code != nil && *op.code != "" {
		query = query.Where("code = ?", *op.code)
	}
	if op.status != nil && *op.status != "" {
		query = query.Where("status = ?", *op.status)
	}
	if op.usedAtStart != nil && *op.usedAtStart != "" {
		query = query.Where("used_at >= ?", *op.usedAtStart)
	}
	if op.usedAtEnd != nil && *op.usedAtEnd != "" {
		query = query.Where("used_at <= ?", *op.usedAtEnd)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := op.page
	pageSize := op.pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 200 {
		pageSize = 200
	}

	offset := (page - 1) * pageSize
	query = query.Order("used_at DESC").Offset(offset).Limit(pageSize)

	var records []model.InvitationRecord
	if err := query.Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}
