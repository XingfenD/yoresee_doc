package notification_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

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
	return op.repo.db.Model(&model.Notification{}).
		Where("id = ? AND (external_id IS NULL OR external_id = '')", op.id).
		Updates(map[string]any{"external_id": op.externalID}).
		Error
}
