package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type NotificationCreateBatchOperation struct {
	repo  *NotificationRepository
	items []model.Notification
	tx    *gorm.DB
}

func (r *NotificationRepository) CreateBatch(items []model.Notification) *NotificationCreateBatchOperation {
	return &NotificationCreateBatchOperation{
		repo:  r,
		items: items,
	}
}

func (op *NotificationCreateBatchOperation) WithTx(tx *gorm.DB) *NotificationCreateBatchOperation {
	op.tx = tx
	return op
}

func (op *NotificationCreateBatchOperation) Exec() error {
	if len(op.items) == 0 {
		return nil
	}
	db := op.repo.db
	if op.tx != nil {
		db = op.tx
	}
	return db.Create(&op.items).Error
}
