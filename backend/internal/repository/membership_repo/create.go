package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type CreateUserGroupOperation struct {
	repo  *MembershipRepository
	group *model.UserGroupMeta
	tx    *gorm.DB
}

func (r *MembershipRepository) CreateUserGroup(group *model.UserGroupMeta) *CreateUserGroupOperation {
	return &CreateUserGroupOperation{
		repo:  r,
		group: group,
	}
}

func (op *CreateUserGroupOperation) WithTx(tx *gorm.DB) *CreateUserGroupOperation {
	op.tx = tx
	return op
}

func (op *CreateUserGroupOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.group).Error
	}
	return storage.DB.Create(op.group).Error
}

type CreateOrgNodeOperation struct {
	repo    *MembershipRepository
	orgNode *model.OrgNodeMeta
	tx      *gorm.DB
}

func (r *MembershipRepository) CreateOrgNode(orgNode *model.OrgNodeMeta) *CreateOrgNodeOperation {
	return &CreateOrgNodeOperation{
		repo:    r,
		orgNode: orgNode,
	}
}

func (op *CreateOrgNodeOperation) WithTx(tx *gorm.DB) *CreateOrgNodeOperation {
	op.tx = tx
	return op
}

func (op *CreateOrgNodeOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.orgNode).Error
	}
	return storage.DB.Create(op.orgNode).Error
}
