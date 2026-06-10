package document_version_repo

import "gorm.io/gorm"

type DocumentVersionRepository struct {
	db *gorm.DB
}

func NewDocumentVersionRepository(db *gorm.DB) *DocumentVersionRepository {
	return &DocumentVersionRepository{db: db}
}
