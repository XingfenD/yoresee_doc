package user_repo

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type QueryUsersOperation struct {
	repo     *UserRepository
	tx       *gorm.DB
	keyword  *string
	page     int
	pageSize int
}

func (r *UserRepository) QueryUsers() *QueryUsersOperation {
	return &QueryUsersOperation{
		repo: r,
	}
}

func (op *QueryUsersOperation) WithTx(tx *gorm.DB) *QueryUsersOperation {
	op.tx = tx
	return op
}

func (op *QueryUsersOperation) WithKeyword(keyword *string) *QueryUsersOperation {
	op.keyword = keyword
	return op
}

func (op *QueryUsersOperation) WithPagination(page, pageSize int) *QueryUsersOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *QueryUsersOperation) ExecWithTotal() ([]model.User, int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(&model.User{})
	if op.keyword != nil {
		trimmed := strings.TrimSpace(*op.keyword)
		if trimmed != "" {
			like := "%" + trimmed + "%"
			query = query.Where("username ILIKE ? OR email ILIKE ? OR external_id ILIKE ?", like, like, like)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		query = query.Offset(offset).Limit(op.pageSize)
	}

	var users []model.User
	if err := query.Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
