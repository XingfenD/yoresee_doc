package knowledge_base_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ListKnowledgeBaseOperation struct {
	repo          *KnowledgeBaseRepository
	model         *model.KnowledgeBase
	creatorID     *int64
	isPublic      *bool
	nameKeyword   *string
	createAtStart *string
	createAtEnd   *string
	updateAtStart *string
	updateAtEnd   *string
	sortField     string
	sortDesc      bool
	page          int
	pageSize      int
	tx            *gorm.DB
}

func (r *KnowledgeBaseRepository) List(m *model.KnowledgeBase) (op *ListKnowledgeBaseOperation) {
	return &ListKnowledgeBaseOperation{
		repo:  KnowledgeBaseRepo,
		model: m,
	}
}

func (op *ListKnowledgeBaseOperation) WithTx(tx *gorm.DB) *ListKnowledgeBaseOperation {
	op.tx = tx
	return op
}

func (op *ListKnowledgeBaseOperation) WithCreatorID(creatorID *int64) *ListKnowledgeBaseOperation {
	op.creatorID = creatorID
	return op
}

func (op *ListKnowledgeBaseOperation) WithIsPublic(isPublic *bool) *ListKnowledgeBaseOperation {
	op.isPublic = isPublic
	return op
}

func (op *ListKnowledgeBaseOperation) WithNameKeyword(nameKeyword *string) *ListKnowledgeBaseOperation {
	op.nameKeyword = nameKeyword
	return op
}

func (op *ListKnowledgeBaseOperation) WithCreateTimeRange(start, end *string) *ListKnowledgeBaseOperation {
	op.createAtStart = start
	op.createAtEnd = end
	return op
}

func (op *ListKnowledgeBaseOperation) WithUpdateTimeRange(start, end *string) *ListKnowledgeBaseOperation {
	op.updateAtStart = start
	op.updateAtEnd = end
	return op
}

func (op *ListKnowledgeBaseOperation) WithSort(field string, desc bool) *ListKnowledgeBaseOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *ListKnowledgeBaseOperation) WithPagination(page, pageSize int) *ListKnowledgeBaseOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *ListKnowledgeBaseOperation) Exec() (kbs []*model.KnowledgeBase, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	dbQuery := op.tx.Model(op.model)

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("creator_user_id = ?", *op.creatorID)
	}

	if op.isPublic != nil {
		dbQuery = dbQuery.Where("is_public = ?", *op.isPublic)
	}

	orderStr := "created_at DESC"
	if op.sortField != "" {
		if op.sortDesc {
			orderStr = op.sortField + " DESC"
		} else {
			orderStr = op.sortField + " ASC"
		}
	}
	dbQuery = dbQuery.Order(orderStr)

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		dbQuery = dbQuery.Offset(offset).Limit(op.pageSize)
	}

	err = dbQuery.Find(&kbs).Error
	return
}

func (op *ListKnowledgeBaseOperation) ExecWithTotal() (kbs []*model.KnowledgeBase, total int64, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	dbQuery := op.tx.Model(op.model)

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.creatorID != nil {
		dbQuery = dbQuery.Where("creator_user_id = ?", *op.creatorID)
	}

	if op.isPublic != nil {
		dbQuery = dbQuery.Where("is_public = ?", *op.isPublic)
	}
	err = dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	orderStr := "created_at DESC"
	if op.sortField != "" {
		if op.sortDesc {
			orderStr = op.sortField + " DESC"
		} else {
			orderStr = op.sortField + " ASC"
		}
	}
	dbQuery = dbQuery.Order(orderStr)

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		dbQuery = dbQuery.Offset(offset).Limit(op.pageSize)
	}

	err = dbQuery.Find(&kbs).Error
	return
}
