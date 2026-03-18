package invitation_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationGetByIDOperation struct {
	repo *InvitationRepository
	id   int64
	tx   *gorm.DB
}

func (r *InvitationRepository) GetByID(id int64) *InvitationGetByIDOperation {
	return &InvitationGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *InvitationGetByIDOperation) WithTx(tx *gorm.DB) *InvitationGetByIDOperation {
	op.tx = tx
	return op
}

func (op *InvitationGetByIDOperation) Exec() (*model.Invitation, error) {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.First(&invitation, op.id).Error
	} else {
		err = storage.DB.First(&invitation, op.id).Error
	}

	return &invitation, err
}

type InvitationCreateOperation struct {
	repo       *InvitationRepository
	invitation *model.Invitation
	tx         *gorm.DB
}

func (r *InvitationRepository) Create(invitation *model.Invitation) *InvitationCreateOperation {
	return &InvitationCreateOperation{
		repo:       r,
		invitation: invitation,
	}
}

func (op *InvitationCreateOperation) WithTx(tx *gorm.DB) *InvitationCreateOperation {
	op.tx = tx
	return op
}

func (op *InvitationCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.invitation).Error
	}
	return storage.DB.Create(op.invitation).Error
}

type InvitationUpdateOperation struct {
	repo       *InvitationRepository
	invitation *model.Invitation
	tx         *gorm.DB
}

func (r *InvitationRepository) Update(invitation *model.Invitation) *InvitationUpdateOperation {
	return &InvitationUpdateOperation{
		repo:       r,
		invitation: invitation,
	}
}

func (op *InvitationUpdateOperation) WithTx(tx *gorm.DB) *InvitationUpdateOperation {
	op.tx = tx
	return op
}

func (op *InvitationUpdateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Save(op.invitation).Error
	}
	return storage.DB.Save(op.invitation).Error
}

type InvitationDeleteOperation struct {
	repo *InvitationRepository
	id   int64
	tx   *gorm.DB
}

func (r *InvitationRepository) Delete(id int64) *InvitationDeleteOperation {
	return &InvitationDeleteOperation{
		repo: r,
		id:   id,
	}
}

func (op *InvitationDeleteOperation) WithTx(tx *gorm.DB) *InvitationDeleteOperation {
	op.tx = tx
	return op
}

func (op *InvitationDeleteOperation) Exec() error {
	if op.tx == nil {
		op.tx = storage.DB
	}
	return op.tx.Delete(&model.Invitation{}, op.id).Error
}
