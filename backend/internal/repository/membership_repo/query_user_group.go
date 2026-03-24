package membership_repo

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type QueryUserGroupOperation struct {
	repo     *MembershipRepository
	tx       *gorm.DB
	keyword  *string
	page     int
	pageSize int
}

func (r *MembershipRepository) QueryUserGroup() *QueryUserGroupOperation {
	return &QueryUserGroupOperation{
		repo: r,
	}
}

func (op *QueryUserGroupOperation) WithTx(tx *gorm.DB) *QueryUserGroupOperation {
	op.tx = tx
	return op
}

func (op *QueryUserGroupOperation) WithKeyword(keyword *string) *QueryUserGroupOperation {
	op.keyword = keyword
	return op
}

func (op *QueryUserGroupOperation) WithPagination(page, pageSize int) *QueryUserGroupOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *QueryUserGroupOperation) ExecWithTotal() ([]model.UserGroupMeta, int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(&model.UserGroupMeta{})
	if op.keyword != nil {
		trimmed := strings.TrimSpace(*op.keyword)
		if trimmed != "" {
			like := "%" + trimmed + "%"
			query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
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

	var userGroups []model.UserGroupMeta
	if err := query.Order("id DESC").Find(&userGroups).Error; err != nil {
		return nil, 0, err
	}
	return userGroups, total, nil
}
