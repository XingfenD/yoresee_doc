package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func (r *DocumentRepository) MGetByIDs(ids []int64) ([]*model.Document, error) {
	if len(ids) == 0 {
		return []*model.Document{}, nil
	}

	var docs []*model.Document
	if err := storage.DB.Model(&model.Document{}).Where("id IN ?", ids).Find(&docs).Error; err != nil {
		return nil, err
	}

	docMap := make(map[int64]*model.Document, len(docs))
	for _, doc := range docs {
		docMap[doc.ID] = doc
	}

	ordered := make([]*model.Document, 0, len(ids))
	for _, id := range ids {
		if doc, ok := docMap[id]; ok {
			ordered = append(ordered, doc)
		}
	}
	return ordered, nil
}
