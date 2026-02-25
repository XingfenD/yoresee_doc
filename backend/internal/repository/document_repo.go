package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentRepository struct{}

var DocumentRepo = &DocumentRepository{}

type DocumentGetByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	tx         *gorm.DB
}

func (r *DocumentRepository) GetByExternalID(externalID string) *DocumentGetByExternalIDOperation {
	return &DocumentGetByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *DocumentGetByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentGetByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetByExternalIDOperation) Exec() (*model.DocumentMeta, error) {
	var document model.DocumentMeta
	var err error

	if op.tx != nil {
		err = op.tx.First(&document, "external_id = ?", op.externalID).Error
	} else {
		err = storage.DB.First(&document, "external_id = ?", op.externalID).Error
	}

	return &document, err
}

type DocumentGetContentOperation struct {
	repo       *DocumentRepository
	documentID int64
	tx         *gorm.DB
}

func (r *DocumentRepository) GetContent(documentID int64) *DocumentGetContentOperation {
	return &DocumentGetContentOperation{
		repo:       r,
		documentID: documentID,
	}
}

func (op *DocumentGetContentOperation) WithTx(tx *gorm.DB) *DocumentGetContentOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetContentOperation) Exec() (string, error) {
	var version model.DocumentVersion
	var content model.Content
	var err error

	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	err = db.Where("document_id = ?", op.documentID).Order("version DESC").First(&version).Error
	if err != nil {
		return "", err
	}

	err = db.First(&content, version.ContentID).Error
	if err != nil {
		return "", err
	}

	return content.Content, nil
}

type ListDocumentsOperation struct {
	repo                 *DocumentRepository
	model                *model.DocumentMeta
	userID               *int64
	parentID             *int64
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
	tx                   *gorm.DB
}

func (r *DocumentRepository) ListDocuments(documentModel *model.DocumentMeta) *ListDocumentsOperation {
	return &ListDocumentsOperation{
		repo:  r,
		model: documentModel,
	}
}

func (op *ListDocumentsOperation) WithUserID(userID *int64) *ListDocumentsOperation {
	op.userID = userID
	return op
}

func (op *ListDocumentsOperation) WithParentID(parentID *int64) *ListDocumentsOperation {
	op.parentID = parentID
	return op
}

func (op *ListDocumentsOperation) WithTitleKeyword(titleKeyword *string) *ListDocumentsOperation {
	op.titleKeyword = titleKeyword
	return op
}

func (op *ListDocumentsOperation) WithType(docType *string) *ListDocumentsOperation {
	op.docType = docType
	return op
}

func (op *ListDocumentsOperation) WithStatus(status *int) *ListDocumentsOperation {
	op.status = status
	return op
}

func (op *ListDocumentsOperation) WithTags(tags []string) *ListDocumentsOperation {
	op.tags = tags
	return op
}

func (op *ListDocumentsOperation) WithCreateTimeRange(start, end *string) *ListDocumentsOperation {
	op.createTimeRangeStart = start
	op.createTimeRangeEnd = end
	return op
}

func (op *ListDocumentsOperation) WithUpdateTimeRange(start, end *string) *ListDocumentsOperation {
	op.updateTimeRangeStart = start
	op.updateTimeRangeEnd = end
	return op
}

func (op *ListDocumentsOperation) WithPagination(page, pageSize int) *ListDocumentsOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *ListDocumentsOperation) WithSort(field string, desc bool) *ListDocumentsOperation {
	op.sortField = field
	op.sortDesc = desc
	return op
}

func (op *ListDocumentsOperation) WithTx(tx *gorm.DB) *ListDocumentsOperation {
	op.tx = tx
	return op
}

func (op *ListDocumentsOperation) Exec() ([]model.DocumentMeta, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	dbQuery := db.Model(&model.DocumentMeta{})

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.userID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *op.userID)
	}

	if op.parentID != nil {
		dbQuery = dbQuery.Where("parent_id = ?", *op.parentID)
	}

	if op.titleKeyword != nil && *op.titleKeyword != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+*op.titleKeyword+"%")
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

	var documents []model.DocumentMeta
	err := dbQuery.Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (op *ListDocumentsOperation) ExecWithTotal() ([]model.DocumentMeta, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	dbQuery := db.Model(&model.DocumentMeta{})

	if op.model != nil {
		dbQuery = dbQuery.Where(op.model)
	}

	if op.userID != nil {
		dbQuery = dbQuery.Where("user_id = ?", *op.userID)
	}

	if op.parentID != nil {
		dbQuery = dbQuery.Where("parent_id = ?", *op.parentID)
	}

	if op.titleKeyword != nil && *op.titleKeyword != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+*op.titleKeyword+"%")
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

	var total int64
	err := dbQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

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

	var documents []model.DocumentMeta
	err = dbQuery.Find(&documents).Error
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}
