package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DeleteUserGroupOperation struct {
	repo  *MembershipRepository
	group *model.UserGroupMeta
	tx    *gorm.DB
}

func (r *MembershipRepository) DeleteUserGroup(group *model.UserGroupMeta) *DeleteUserGroupOperation {
	return &DeleteUserGroupOperation{
		repo:  r,
		group: group,
	}
}

func (op *DeleteUserGroupOperation) WithTx(tx *gorm.DB) *DeleteUserGroupOperation {
	op.tx = tx
	return op
}

func (op *DeleteUserGroupOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Delete(op.group).Error
	}
	return storage.DB.Delete(op.group).Error
}

type DeleteOrgNodeOperation struct {
	repo    *MembershipRepository
	orgNode *model.OrgNodeMeta
	tx      *gorm.DB
}

func (r *MembershipRepository) DeleteOrgNode(orgNode *model.OrgNodeMeta) *DeleteOrgNodeOperation {
	return &DeleteOrgNodeOperation{
		repo:    r,
		orgNode: orgNode,
	}
}

func (op *DeleteOrgNodeOperation) WithTx(tx *gorm.DB) *DeleteOrgNodeOperation {
	op.tx = tx
	return op
}

func (op *DeleteOrgNodeOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Delete(op.orgNode).Error
	}
	return storage.DB.Delete(op.orgNode).Error
}
