package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentUpdateOperation struct {
	repo         *DocumentRepository
	doc          *model.Document
	updateFields map[string]bool
	tx           *gorm.DB
}

func (r *DocumentRepository) Update(doc *model.Document) *DocumentUpdateOperation {
	return &DocumentUpdateOperation{
		repo:         r,
		doc:          doc,
		updateFields: make(map[string]bool),
	}
}

func (op *DocumentUpdateOperation) UpdateTitle() *DocumentUpdateOperation {
	op.updateFields["title"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateSummary() *DocumentUpdateOperation {
	op.updateFields["summary"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateContent() *DocumentUpdateOperation {
	op.updateFields["content"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateParentID() *DocumentUpdateOperation {
	op.updateFields["parent_id"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateKnowledgeID() *DocumentUpdateOperation {
	op.updateFields["knowledge_id"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateContainerType() *DocumentUpdateOperation {
	op.updateFields["container_type"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateIsPublic() *DocumentUpdateOperation {
	op.updateFields["is_public"] = true
	return op
}

func (op *DocumentUpdateOperation) UpdateTags() *DocumentUpdateOperation {
	op.updateFields["tags"] = true
	return op
}

func (op *DocumentUpdateOperation) WithTx(tx *gorm.DB) *DocumentUpdateOperation {
	op.tx = tx
	return op
}

func (op *DocumentUpdateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(op.doc)

	if len(op.updateFields) > 0 {
		fields := make([]string, 0, len(op.updateFields))
		for field := range op.updateFields {
			fields = append(fields, field)
		}
		query = query.Select(fields)
	}

	return query.Updates(*op.doc).Error
}
