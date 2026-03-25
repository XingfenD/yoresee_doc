package membership_repo

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type QueryOrgNodeOperation struct {
	repo     *MembershipRepository
	tx       *gorm.DB
	parentID int64
	keyword  *string
	page     int
	pageSize int
}

func (r *MembershipRepository) QueryOrgNode() *QueryOrgNodeOperation {
	return &QueryOrgNodeOperation{
		repo: r,
	}
}

func (op *QueryOrgNodeOperation) WithTx(tx *gorm.DB) *QueryOrgNodeOperation {
	op.tx = tx
	return op
}

func (op *QueryOrgNodeOperation) WithParentID(parentID int64) *QueryOrgNodeOperation {
	op.parentID = parentID
	return op
}

func (op *QueryOrgNodeOperation) WithKeyword(keyword *string) *QueryOrgNodeOperation {
	op.keyword = keyword
	return op
}

func (op *QueryOrgNodeOperation) WithPagination(page, pageSize int) *QueryOrgNodeOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *QueryOrgNodeOperation) ExecWithTotal() ([]model.OrgNodeMeta, int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(&model.OrgNodeMeta{})
	if op.parentID >= 0 {
		query = query.Where("parent_id = ?", op.parentID)
	}
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

	var orgNodes []model.OrgNodeMeta
	if err := query.Order("id DESC").Find(&orgNodes).Error; err != nil {
		return nil, 0, err
	}
	return orgNodes, total, nil
}

func (op *QueryOrgNodeOperation) Exec() ([]model.OrgNodeMeta, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	query := op.tx.Model(&model.OrgNodeMeta{})
	if op.parentID >= 0 {
		query = query.Where("parent_id = ?", op.parentID)
	}
	if op.keyword != nil {
		trimmed := strings.TrimSpace(*op.keyword)
		if trimmed != "" {
			like := "%" + trimmed + "%"
			query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
		}
	}

	var orgNodes []model.OrgNodeMeta
	if err := query.Order("id DESC").Find(&orgNodes).Error; err != nil {
		return nil, err
	}
	return orgNodes, nil
}

type MGetOrgNodeByIDOperation struct {
	repo    *MembershipRepository
	nodeIDs []int64
	tx      *gorm.DB
}

func (r *MembershipRepository) MGetOrgNodeByID(nodeIDs []int64) *MGetOrgNodeByIDOperation {
	return &MGetOrgNodeByIDOperation{
		repo:    r,
		nodeIDs: nodeIDs,
	}
}

func (op *MGetOrgNodeByIDOperation) WithTx(tx *gorm.DB) *MGetOrgNodeByIDOperation {
	op.tx = tx
	return op
}

func (op *MGetOrgNodeByIDOperation) Exec() (map[int64]*model.OrgNodeMeta, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	var orgNodes []model.OrgNodeMeta
	if err := op.tx.Where("id IN ?", op.nodeIDs).Find(&orgNodes).Error; err != nil {
		return nil, err
	}

	result := make(map[int64]*model.OrgNodeMeta, len(orgNodes))
	for i := range orgNodes {
		result[orgNodes[i].ID] = &orgNodes[i]
	}
	return result, nil
}

type QueryOrgNodeByPathPrefixOperation struct {
	repo   *MembershipRepository
	tx     *gorm.DB
	prefix string
}

func (r *MembershipRepository) QueryOrgNodeByPathPrefix(prefix string) *QueryOrgNodeByPathPrefixOperation {
	return &QueryOrgNodeByPathPrefixOperation{
		repo:   r,
		prefix: prefix,
	}
}

func (op *QueryOrgNodeByPathPrefixOperation) WithTx(tx *gorm.DB) *QueryOrgNodeByPathPrefixOperation {
	op.tx = tx
	return op
}

func (op *QueryOrgNodeByPathPrefixOperation) Exec() ([]model.OrgNodeMeta, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	var orgNodes []model.OrgNodeMeta
	prefix := strings.TrimSpace(op.prefix)
	if prefix == "" {
		return []model.OrgNodeMeta{}, nil
	}
	like := prefix + ".%"
	if err := op.tx.Where("path = ? OR path LIKE ?", prefix, like).Find(&orgNodes).Error; err != nil {
		return nil, err
	}
	return orgNodes, nil
}
