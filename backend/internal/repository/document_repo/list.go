package document_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentsListOperation struct {
	repo                 *DocumentRepository
	model                *model.Document
	userID               *int64
	parentID             *int64
	knowledgeID          *int64
	listOwnDoc           bool
	ids                  []int64
	titleKeyword         *string
	docType              *string
	status               *int
	tags                 []string
	createTimeRangeStart *string
	createTimeRangeEnd   *string
	updateTimeRangeStart *string
	updateTimeRangeEnd   *string
	sortField            string
	sortDesc             bool
	page                 int
	pageSize             int

	directoryOnly bool

	tx *gorm.DB
}

func (r *DocumentRepository) ListDocuments(documentModel *model.Document) *DocumentsListOperation {
	return &DocumentsListOperation{
		repo:  r,
		model: documentModel,
	}
}

func (op *DocumentsListOperation) WithUserID(userID *int64) *DocumentsListOperation {
	op.userID = userID
	return op
}

func (op *DocumentsListOperation) WithParentID(parentID *int64) *DocumentsListOperation {
	op.parentID = parentID
	return op
}

func (op *DocumentsListOperation) WithKnowledgeID(knowledgeID *int64) *DocumentsListOperation {
	op.knowledgeID = knowledgeID
	return op
}

func (op *DocumentsListOperation) WithListOwnDoc(listOwnDoc bool) *DocumentsListOperation {
	op.listOwnDoc = listOwnDoc
	return op
}

func (op *DocumentsListOperation) WithIDs(ids []int64) *DocumentsListOperation {
	op.ids = ids
	return op
}

func (op *DocumentsListOperation) WithTitleKeyword(titleKeyword *string) *DocumentsListOperation {
	op.titleKeyword = titleKeyword
	return op
}

func (op *DocumentsListOperation) WithType(docType *string) *DocumentsListOperation {
	op.docType = docType
	return op
}

func (op *DocumentsListOperation) WithStatus(status *int) *DocumentsListOperation {
	op.status = status
	return op
}

func (op *DocumentsListOperation) WithTags(tags []string) *DocumentsListOperation {
	op.tags = tags
	return op
}

func (op *DocumentsListOperation) WithCreateTimeRange(start, end *string) *DocumentsListOperation {
	op.createTimeRangeStart = start
	op.createTimeRangeEnd = end
	return op
}

func (op *DocumentsListOperation) WithUpdateTimeRange(start, end *string) *DocumentsListOperation {
	op.updateTimeRangeStart = start
	op.updateTimeRangeEnd = end
	return op
}

func (op *DocumentsListOperation) WithPagination(page, pageSize int) *DocumentsListOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *DocumentsListOperation) WithSort(field string, desc bool) *DocumentsListOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *DocumentsListOperation) WithDirectoryOnly(with bool) *DocumentsListOperation {
	op.directoryOnly = with
	return op
}

func (op *DocumentsListOperation) WithTx(tx *gorm.DB) *DocumentsListOperation {
	op.tx = tx
	return op
}

func (op *DocumentsListOperation) buildBaseQuery() *gorm.DB {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	dbQuery := db.Model(&model.Document{})
	if op.directoryOnly {
		dbQuery = dbQuery.Select("id, external_id, title, parent_id")
	}

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.userID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *op.userID)
	}

	if op.parentID != nil {
		dbQuery = dbQuery.Where("parent_id = ?", *op.parentID)
	}

	if op.listOwnDoc {
		dbQuery = dbQuery.Where("knowledge_id IS NULL")
	}

	if op.ids != nil {
		if len(op.ids) == 0 {
			dbQuery = dbQuery.Where("1 = 0")
		} else {
			dbQuery = dbQuery.Where("id IN ?", op.ids)
		}
	}

	if op.knowledgeID != nil {
		dbQuery = dbQuery.Where("knowledge_id = ?", *op.knowledgeID)
	}

	if op.titleKeyword != nil && *op.titleKeyword != "" {
		like := "%" + *op.titleKeyword + "%"
		dbQuery = dbQuery.Where("(title LIKE ? OR content LIKE ?)", like, like)
	}

	if op.docType != nil && *op.docType != "" {
		dbQuery = dbQuery.Where("type = ?", *op.docType)
	}

	if op.status != nil {
		dbQuery = dbQuery.Where("status = ?", *op.status)
	}

	if len(op.tags) > 0 {
		for _, tag := range op.tags {
			dbQuery = dbQuery.Where("JSON_CONTAINS(tags, ?)", "\""+tag+"\"")
		}
	}

	if op.createTimeRangeStart != nil && *op.createTimeRangeStart != "" {
		dbQuery = dbQuery.Where("created_at >= ?", *op.createTimeRangeStart)
	}
	if op.createTimeRangeEnd != nil && *op.createTimeRangeEnd != "" {
		dbQuery = dbQuery.Where("created_at <= ?", *op.createTimeRangeEnd)
	}
	if op.updateTimeRangeStart != nil && *op.updateTimeRangeStart != "" {
		dbQuery = dbQuery.Where("updated_at >= ?", *op.updateTimeRangeStart)
	}
	if op.updateTimeRangeEnd != nil && *op.updateTimeRangeEnd != "" {
		dbQuery = dbQuery.Where("updated_at <= ?", *op.updateTimeRangeEnd)
	}

	return dbQuery
}

func (op *DocumentsListOperation) appendOtherArgs(db *gorm.DB) *gorm.DB {
	dbQuery := db

	sortField := op.sortField
	if sortField == "" {
		sortField = "created_at"
	}
	orderDirection := "ASC"
	if op.sortDesc {
		orderDirection = "DESC"
	}
	dbQuery = dbQuery.Order(sortField + " " + orderDirection)

	page := op.page
	pageSize := op.pageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	dbQuery = dbQuery.Offset(offset).Limit(pageSize)

	return dbQuery
}

func (op *DocumentsListOperation) Exec() ([]model.Document, error) {
	dbQuery := op.buildBaseQuery()
	dbQuery = op.appendOtherArgs(dbQuery)

	var documents []model.Document
	err := dbQuery.Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (op *DocumentsListOperation) ExecWithTotal() ([]model.Document, int64, error) {
	dbQuery := op.buildBaseQuery()

	var total int64
	err := dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	dbQuery = op.appendOtherArgs(dbQuery)

	var documents []model.Document
	err = dbQuery.Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}
