package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type ListUserGroupOperation struct {
	repo *MembershipRepository
	tx   *gorm.DB
}

func (r *MembershipRepository) ListUserGroup() *ListUserGroupOperation {
	return &ListUserGroupOperation{
		repo: r,
	}
}

func (op *ListUserGroupOperation) WithTx(tx *gorm.DB) *ListUserGroupOperation {
	op.tx = tx
	return op
}

func (op *ListUserGroupOperation) Exec() ([]model.UserGroupMeta, error) {
	var userGroups []model.UserGroupMeta
	var err error

	if op.tx != nil {
		err = op.tx.Find(&userGroups).Error
	} else {
		err = op.repo.db.Find(&userGroups).Error
	}

	return userGroups, err
}
