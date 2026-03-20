package template_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ListTemplateOperation struct {
	repo            *TemplateRepository
	model           *model.Template
	userID          *int64
	scope           *string
	knowledgeBaseID *int64
	nameKeyword     *string
	sortField       string
	sortDesc        bool
	page            int
	pageSize        int
	tx              *gorm.DB
}

func (r *TemplateRepository) List(m *model.Template) (op *ListTemplateOperation) {
	return &ListTemplateOperation{
		repo:  TemplateRepo,
		model: m,
	}
}

func (op *ListTemplateOperation) WithTx(tx *gorm.DB) *ListTemplateOperation {
	op.tx = tx
	return op
}

func (op *ListTemplateOperation) WithUserID(userID *int64) *ListTemplateOperation {
	op.userID = userID
	return op
}

func (op *ListTemplateOperation) WithScope(scope *string) *ListTemplateOperation {
	op.scope = scope
	return op
}

func (op *ListTemplateOperation) WithKnowledgeBaseID(knowledgeBaseID *int64) *ListTemplateOperation {
	op.knowledgeBaseID = knowledgeBaseID
	return op
}

func (op *ListTemplateOperation) WithNameKeyword(nameKeyword *string) *ListTemplateOperation {
	op.nameKeyword = nameKeyword
	return op
}

func (op *ListTemplateOperation) WithSort(field string, desc bool) *ListTemplateOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *ListTemplateOperation) WithPagination(page, pageSize int) *ListTemplateOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *ListTemplateOperation) ExecWithTotal() (templates []*model.Template, total int64, err error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	dbQuery := op.tx.Model(op.model)
	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.userID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *op.userID)
	}
	if op.scope != nil {
		dbQuery = dbQuery.Where("scope = ?", *op.scope)
	}
	if op.knowledgeBaseID != nil {
		dbQuery = dbQuery.Where("knowledge_base_id = ?", *op.knowledgeBaseID)
	}
	if op.nameKeyword != nil && *op.nameKeyword != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+*op.nameKeyword+"%")
	}

	if err = dbQuery.Count(&total).Error; err != nil {
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

	err = dbQuery.Find(&templates).Error
	return
}
