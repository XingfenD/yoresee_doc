package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type NotificationMarkReadOperation struct {
	repo        *NotificationRepository
	receiverID  int64
	externalIDs []string
	tx          *gorm.DB
}

func (r *NotificationRepository) MarkRead(receiverID int64, externalIDs []string) *NotificationMarkReadOperation {
	return &NotificationMarkReadOperation{
		repo:        r,
		receiverID:  receiverID,
		externalIDs: externalIDs,
	}
}

func (op *NotificationMarkReadOperation) WithTx(tx *gorm.DB) *NotificationMarkReadOperation {
	op.tx = tx
	return op
}

func (op *NotificationMarkReadOperation) Exec() error {
	if len(op.externalIDs) == 0 {
		return nil
	}
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Model(&model.Notification{}).
		Where("receiver_id = ? AND external_id IN ?", op.receiverID, op.externalIDs).
		Updates(map[string]any{"status": "read"}).
		Error
}

type NotificationMarkAllReadOperation struct {
	repo       *NotificationRepository
	receiverID int64
	tx         *gorm.DB
}

func (r *NotificationRepository) MarkAllRead(receiverID int64) *NotificationMarkAllReadOperation {
	return &NotificationMarkAllReadOperation{
		repo:       r,
		receiverID: receiverID,
	}
}

func (op *NotificationMarkAllReadOperation) WithTx(tx *gorm.DB) *NotificationMarkAllReadOperation {
	op.tx = tx
	return op
}

func (op *NotificationMarkAllReadOperation) Exec() error {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}
	return db.Model(&model.Notification{}).
		Where("receiver_id = ? AND status != ?", op.receiverID, "read").
		Updates(map[string]any{"status": "read"}).
		Error
}
