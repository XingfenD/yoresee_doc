package document_version_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func (r *DocumentVersionRepository) GetByDocumentIDAndVersion(documentID int64, version int) (*model.DocumentVersion, error) {
	item := &model.DocumentVersion{}
	if err := storage.DB.Where("document_id = ? AND version = ?", documentID, version).First(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
