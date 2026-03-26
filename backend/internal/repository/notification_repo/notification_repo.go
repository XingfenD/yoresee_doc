package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type NotificationRepository struct{}

var NotificationRepo = &NotificationRepository{}

type NotificationUpdateExternalIDOperation struct {
	repo       *NotificationRepository
	id         int64
	externalID string
}

func (r *NotificationRepository) UpdateExternalID(id int64, externalID string) *NotificationUpdateExternalIDOperation {
	return &NotificationUpdateExternalIDOperation{
		repo:       r,
		id:         id,
		externalID: externalID,
	}
}

func (op *NotificationUpdateExternalIDOperation) Exec() error {
	if op.id == 0 || op.externalID == "" {
		return nil
	}
	return storage.DB.Model(&model.Notification{}).
		Where("id = ? AND (external_id IS NULL OR external_id = '')", op.id).
		Updates(map[string]any{"external_id": op.externalID}).
		Error
}
