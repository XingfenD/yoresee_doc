package document_version_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func (r *DocumentVersionRepository) ListByDocumentID(documentID int64, offset, limit int) ([]*model.DocumentVersion, int64, error) {
	var total int64
	query := storage.DB.Model(&model.DocumentVersion{}).Where("document_id = ?", documentID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*model.DocumentVersion, 0)
	if err := query.Order("version DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
