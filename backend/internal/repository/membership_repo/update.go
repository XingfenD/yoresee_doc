package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UpdateUserGroupOperation struct {
	repo  *MembershipRepository
	group *model.UserGroupMeta
	tx    *gorm.DB
}

func (r *MembershipRepository) UpdateUserGroup(group *model.UserGroupMeta) *UpdateUserGroupOperation {
	return &UpdateUserGroupOperation{
		repo:  r,
		group: group,
	}
}

func (op *UpdateUserGroupOperation) WithTx(tx *gorm.DB) *UpdateUserGroupOperation {
	op.tx = tx
	return op
}

func (op *UpdateUserGroupOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Save(op.group).Error
	}
	return storage.DB.Save(op.group).Error
}

type UpdateOrgNodeOperation struct {
	repo    *MembershipRepository
	orgNode *model.OrgNodeMeta
	tx      *gorm.DB
}

func (r *MembershipRepository) UpdateOrgNode(orgNode *model.OrgNodeMeta) *UpdateOrgNodeOperation {
	return &UpdateOrgNodeOperation{
		repo:    r,
		orgNode: orgNode,
	}
}

func (op *UpdateOrgNodeOperation) WithTx(tx *gorm.DB) *UpdateOrgNodeOperation {
	op.tx = tx
	return op
}

func (op *UpdateOrgNodeOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Save(op.orgNode).Error
	}
	return storage.DB.Save(op.orgNode).Error
}
