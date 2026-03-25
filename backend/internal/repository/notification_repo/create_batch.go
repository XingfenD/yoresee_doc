package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
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
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Create(&op.items).Error
}
