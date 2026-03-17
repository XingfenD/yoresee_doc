package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type GetUserGroupByExternalIDOperation struct {
	repo       *MembershipRepository
	externalID string
	tx         *gorm.DB
}

func (r *MembershipRepository) GetUserGroupByExternalID(externalID string) *GetUserGroupByExternalIDOperation {
	return &GetUserGroupByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *GetUserGroupByExternalIDOperation) WithTx(tx *gorm.DB) *GetUserGroupByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *GetUserGroupByExternalIDOperation) Exec() (*model.UserGroupMeta, error) {
	var group model.UserGroupMeta
	var err error

	if op.tx != nil {
		err = op.tx.Where("external_id = ?", op.externalID).First(&group).Error
	} else {
		err = storage.DB.Where("external_id = ?", op.externalID).First(&group).Error
	}

	return &group, err
}

type GetOrgNodeByExternalID struct {
	repo       *MembershipRepository
	externalID string
	tx         *gorm.DB
}

func (r *MembershipRepository) GetOrgNodeByExternalID(externalID string) *GetOrgNodeByExternalID {
	return &GetOrgNodeByExternalID{
		repo:       r,
		externalID: externalID,
	}
}

func (op *GetOrgNodeByExternalID) WithTx(tx *gorm.DB) *GetOrgNodeByExternalID {
	op.tx = tx
	return op
}

func (op *GetOrgNodeByExternalID) Exec() (*model.OrgNodeMeta, error) {
	var orgNode model.OrgNodeMeta
	var err error

	if op.tx != nil {
		err = op.tx.Where("external_id = ?", op.externalID).First(&orgNode).Error
	} else {
		err = storage.DB.Where("external_id = ?", op.externalID).First(&orgNode).Error
	}

	return &orgNode, err
}
