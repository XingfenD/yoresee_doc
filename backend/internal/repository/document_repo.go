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

func (op *DocumentGetByExternalIDOperation) Exec() (*model.Document, error) {
	var document model.Document
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
	var docMeta model.Document
	var err error

	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	err = db.Where("id = ?", op.documentID).First(&docMeta).Error
	if err != nil {
		return "", err
	}

	return docMeta.Content, nil
}

type DocumentGetIDByExternalIDOperation struct {
	repo       *DocumentRepository
	externalID string
	tx         *gorm.DB
}

func (r *DocumentRepository) GetIDByExternalID(externalID string) *DocumentGetIDByExternalIDOperation {
	return &DocumentGetIDByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *DocumentGetIDByExternalIDOperation) WithTx(tx *gorm.DB) *DocumentGetIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op DocumentGetIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	if op.tx == nil {
		op.tx = storage.DB
	}

	err := op.tx.Model(&model.Document{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

type DocumentsListOperation struct {
	repo                 *DocumentRepository
	model                *model.Document
	userID               *int64
	parentID             *int64
	knowledgeID          *int64
	listOwnDoc           bool
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

	if op.knowledgeID != nil {
		dbQuery = dbQuery.Where("knowledge_id = ?", *op.knowledgeID)
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

type DocumentGetSubtreeOperation struct {
	repo         *DocumentRepository
	rootParentID int64
	knowledgeID  *int64
	depth        *int
	tx           *gorm.DB
}

func (r *DocumentRepository) GetSubtree(rootParentID int64) *DocumentGetSubtreeOperation {
	return &DocumentGetSubtreeOperation{
		repo:         r,
		rootParentID: rootParentID,
	}
}

func (op *DocumentGetSubtreeOperation) WithTx(tx *gorm.DB) *DocumentGetSubtreeOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetSubtreeOperation) WithKnowledgeID(knowledgeID *int64) *DocumentGetSubtreeOperation {
	op.knowledgeID = knowledgeID
	return op
}

func (op *DocumentGetSubtreeOperation) WithDepth(depth int) *DocumentGetSubtreeOperation {
	op.depth = &depth
	return op
}

func (op *DocumentGetSubtreeOperation) Exec() ([]model.Document, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	var documents []model.Document

	if op.depth != nil && *op.depth == 0 {
		return documents, nil
	}

	depthFilter := ""
	if op.depth != nil {
		depthFilter = "AND depth <= " + string(rune(*op.depth+'0'))
	}

	query := `
		WITH RECURSIVE subtree AS (
			SELECT d.*, 0 as depth
			FROM document_metas d
			WHERE d.parent_id = ? AND d.deleted_at IS NULL
			UNION ALL
			SELECT d.*, s.depth + 1 as depth
			FROM document_metas d
			INNER JOIN subtree s ON d.parent_id = s.id
			WHERE d.deleted_at IS NULL
		)
		SELECT * FROM subtree s
		WHERE 1=1 ` + depthFilter + `
		ORDER BY depth, created_at
	`

	err := db.Raw(query, op.rootParentID).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

type DocumentGetSubtreeByKnowledgeIDOperation struct {
	repo        *DocumentRepository
	knowledgeID int64
	depth       *int
	tx          *gorm.DB
}

func (r *DocumentRepository) GetSubtreeByKnowledgeID(knowledgeID int64) *DocumentGetSubtreeByKnowledgeIDOperation {
	return &DocumentGetSubtreeByKnowledgeIDOperation{
		repo:        r,
		knowledgeID: knowledgeID,
	}
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) WithTx(tx *gorm.DB) *DocumentGetSubtreeByKnowledgeIDOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) WithDepth(depth int) *DocumentGetSubtreeByKnowledgeIDOperation {
	op.depth = &depth
	return op
}

func (op *DocumentGetSubtreeByKnowledgeIDOperation) Exec() ([]model.Document, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	var documents []model.Document

	if op.depth != nil && *op.depth == 0 {
		return documents, nil
	}

	depthFilter := ""
	if op.depth != nil {
		depthFilter = "AND depth <= " + string(rune(*op.depth+'0'))
	}

	query := `
		WITH RECURSIVE subtree AS (
			SELECT d.*, 0 as depth
			FROM document_metas d
			WHERE d.knowledge_id = ? AND d.parent_id = 0 AND d.deleted_at IS NULL
			UNION ALL
			SELECT d.*, s.depth + 1 as depth
			FROM document_metas d
			INNER JOIN subtree s ON d.parent_id = s.id
			WHERE d.deleted_at IS NULL
		)
		SELECT * FROM subtree s
		WHERE 1=1 ` + depthFilter + `
		ORDER BY depth, created_at
	`

	err := db.Raw(query, op.knowledgeID).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	return documents, nil
}

type DocumentCreateOperation struct {
	repo *DocumentRepository
	doc  *model.Document
	tx   *gorm.DB
}

func (r *DocumentRepository) Create(doc *model.Document) *DocumentCreateOperation {
	return &DocumentCreateOperation{
		repo: r,
		doc:  doc,
	}
}

func (op *DocumentCreateOperation) WithTx(tx *gorm.DB) *DocumentCreateOperation {
	op.tx = tx
	return op
}

func (op *DocumentCreateOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}

	return op.tx.Create(op.doc).Error
}

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

func (op *DocumentUpdateOperation) UpdateStatus() *DocumentUpdateOperation {
	op.updateFields["status"] = true
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
