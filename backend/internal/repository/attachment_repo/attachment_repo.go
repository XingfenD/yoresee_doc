package attachment_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) Create(attachment *model.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *AttachmentRepository) ListByDocumentID(documentID int64) ([]*model.Attachment, error) {
	attachments := make([]*model.Attachment, 0)
	if err := r.db.
		Where("document_id = ?", documentID).
		Order("created_at desc").
		Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

func (r *AttachmentRepository) GetByExternalIDAndDocumentID(externalID string, documentID int64) (*model.Attachment, error) {
	attachment := &model.Attachment{}
	if err := r.db.
		Where("external_id = ? AND document_id = ?", externalID, documentID).
		First(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

func (r *AttachmentRepository) DeleteByID(id int64) error {
	return r.db.Delete(&model.Attachment{}, id).Error
}

func IsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
