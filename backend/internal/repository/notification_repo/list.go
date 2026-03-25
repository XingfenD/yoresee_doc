package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type NotificationListOperation struct {
	repo       *NotificationRepository
	receiverID int64
	status     *string
	page       int
	pageSize   int
	tx         *gorm.DB
}

func (r *NotificationRepository) List(receiverID int64) *NotificationListOperation {
	return &NotificationListOperation{
		repo:       r,
		receiverID: receiverID,
	}
}

func (op *NotificationListOperation) WithTx(tx *gorm.DB) *NotificationListOperation {
	op.tx = tx
	return op
}

func (op *NotificationListOperation) WithStatus(status *string) *NotificationListOperation {
	op.status = status
	return op
}

func (op *NotificationListOperation) WithPagination(page, pageSize int) *NotificationListOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *NotificationListOperation) ExecWithTotal() ([]model.Notification, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	query := db.Model(&model.Notification{}).Where("receiver_id = ?", op.receiverID)
	if op.status != nil && *op.status != "" {
		query = query.Where("status = ?", *op.status)
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
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var list []model.Notification
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
