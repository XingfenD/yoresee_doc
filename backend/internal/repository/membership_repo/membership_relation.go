package membership_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
