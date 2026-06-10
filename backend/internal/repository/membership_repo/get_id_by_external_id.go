package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type GetUserGroupIDByExternalIDOperation struct {
	repo       *MembershipRepository
	externalID string
	tx         *gorm.DB
}

func (r *MembershipRepository) GetUserGroupIDByExternalID(externalID string) *GetUserGroupIDByExternalIDOperation {
	return &GetUserGroupIDByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *GetUserGroupIDByExternalIDOperation) WithTx(tx *gorm.DB) *GetUserGroupIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *GetUserGroupIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	var err error
	if op.tx != nil {
		err = op.tx.Model(&model.UserGroupMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	}
	err = op.repo.db.Model(&model.UserGroupMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	return id, err
}

type GetOrgNodeIDByExternalIDOperation struct {
	repo       *MembershipRepository
	externalID string
	tx         *gorm.DB
}

func (r *MembershipRepository) GetOrgNodeIDByExternalID(externalID string) *GetOrgNodeIDByExternalIDOperation {
	return &GetOrgNodeIDByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *GetOrgNodeIDByExternalIDOperation) WithTx(tx *gorm.DB) *GetOrgNodeIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *GetOrgNodeIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	var err error
	if op.tx != nil {
		err = op.tx.Model(&model.OrgNodeMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	}
	err = op.repo.db.Model(&model.OrgNodeMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	return id, err
}
