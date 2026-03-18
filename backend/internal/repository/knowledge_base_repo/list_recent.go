package knowledge_base_repo

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ListRecentKnowledgeBasesOperation struct {
	repo      *KnowledgeBaseRepository
	userID    int64
	startTime *time.Time
	endTime   *time.Time
	page      int
	pageSize  int
	tx        *gorm.DB
}

func (r *KnowledgeBaseRepository) ListRecentKnowledgeBases(userID int64) *ListRecentKnowledgeBasesOperation {
	return &ListRecentKnowledgeBasesOperation{
		repo:   r,
		userID: userID,
	}
}

func (op *ListRecentKnowledgeBasesOperation) WithTimeRange(start, end *time.Time) *ListRecentKnowledgeBasesOperation {
	op.startTime = start
	op.endTime = end
	return op
}

func (op *ListRecentKnowledgeBasesOperation) WithPagination(page, pageSize int) *ListRecentKnowledgeBasesOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *ListRecentKnowledgeBasesOperation) WithTx(tx *gorm.DB) *ListRecentKnowledgeBasesOperation {
	op.tx = tx
	return op
}

func (op *ListRecentKnowledgeBasesOperation) Exec() ([]*model.RecentKnowledgeBase, int64, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	query := db.Model(&model.RecentKnowledgeBase{}).
		Where("user_id = ?", op.userID)
	if op.startTime != nil {
		query = query.Where("accessed_at >= ?", *op.startTime)
	}
	if op.endTime != nil {
		query = query.Where("accessed_at <= ?", *op.endTime)
	}

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

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var records []*model.RecentKnowledgeBase
	if err := query.Order("accessed_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}
