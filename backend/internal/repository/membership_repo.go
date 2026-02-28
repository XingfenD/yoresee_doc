package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type MembershipRepository struct{}

var MembershipRepo = &MembershipRepository{}

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
		err = storage.DB.Find(&userGroups).Error
	}

	return userGroups, err
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

type CreateMembershipOperation struct {
	repo       *MembershipRepository
	membership *model.MembershipRelation
	tx         *gorm.DB
}

func (r *MembershipRepository) CreateMembership(membership *model.MembershipRelation) *CreateMembershipOperation {
	return &CreateMembershipOperation{
		repo:       r,
		membership: membership,
	}
}

func (op *CreateMembershipOperation) WithTx(tx *gorm.DB) *CreateMembershipOperation {
	op.tx = tx
	return op
}

func (op *CreateMembershipOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.membership).Error
	}
	return storage.DB.Create(op.membership).Error
}

type BatchCreateMembershipOperation struct {
	repo        *MembershipRepository
	memberships []*model.MembershipRelation
	tx          *gorm.DB
}

func (r *MembershipRepository) BatchCreateMembership(memberships []*model.MembershipRelation) *BatchCreateMembershipOperation {
	return &BatchCreateMembershipOperation{
		repo:        r,
		memberships: memberships,
	}
}

func (op *BatchCreateMembershipOperation) WithTx(tx *gorm.DB) *BatchCreateMembershipOperation {
	op.tx = tx
	return op
}

func (op *BatchCreateMembershipOperation) Exec() error {
	if len(op.memberships) == 0 {
		return nil
	}

	if op.tx != nil {
		return op.tx.Create(&op.memberships).Error
	}
	return storage.DB.Create(&op.memberships).Error
}

type DeleteMembershipOperation struct {
	repo       *MembershipRepository
	membership *model.MembershipRelation
	tx         *gorm.DB
}

func (r *MembershipRepository) DeleteMembership(membership *model.MembershipRelation) *DeleteMembershipOperation {
	return &DeleteMembershipOperation{
		repo:       r,
		membership: membership,
	}
}

func (op *DeleteMembershipOperation) WithTx(tx *gorm.DB) *DeleteMembershipOperation {
	op.tx = tx
	return op
}

func (op *DeleteMembershipOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Delete(op.membership).Error
	}
	return storage.DB.Delete(op.membership).Error
}

type BatchDeleteMembershipOperation struct {
	repo        *MembershipRepository
	memberships []*model.MembershipRelation
	tx          *gorm.DB
}

func (r *MembershipRepository) BatchDeleteMembership(memberships []*model.MembershipRelation) *BatchDeleteMembershipOperation {
	return &BatchDeleteMembershipOperation{
		repo:        r,
		memberships: memberships,
	}
}

func (op *BatchDeleteMembershipOperation) WithTx(tx *gorm.DB) *BatchDeleteMembershipOperation {
	op.tx = tx
	return op
}

func (op *BatchDeleteMembershipOperation) Exec() error {
	if len(op.memberships) == 0 {
		return nil
	}

	if op.tx != nil {
		return op.tx.Delete(&op.memberships).Error
	}
	return storage.DB.Delete(&op.memberships).Error
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
	err = storage.DB.Model(&model.OrgNodeMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	return id, err
}

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
	err = storage.DB.Model(&model.UserGroupMeta{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	return id, err
}

type ListMembershipOperation struct {
	repo  *MembershipRepository
	query *model.MembershipRelation
	tx    *gorm.DB
}

func (r *MembershipRepository) ListMembership(query *model.MembershipRelation) *ListMembershipOperation {
	return &ListMembershipOperation{
		repo:  r,
		query: query,
	}
}

func (op *ListMembershipOperation) WithTx(tx *gorm.DB) *ListMembershipOperation {
	op.tx = tx
	return op
}

func (op *ListMembershipOperation) Exec() ([]model.MembershipRelation, error) {
	var memberships []model.MembershipRelation
	var err error

	if op.tx != nil {
		err = op.tx.Where(op.query).Find(&memberships).Error
	} else {
		err = storage.DB.Where(op.query).Find(&memberships).Error
	}

	return memberships, err
}
