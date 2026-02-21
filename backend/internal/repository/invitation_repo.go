package repository

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type InvitationRepository struct{}

var InvitationRepo = &InvitationRepository{}

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

type InvitationGetByCodeOperation struct {
	repo *InvitationRepository
	code string
	tx   *gorm.DB
}

func (r *InvitationRepository) GetByCode(code string) *InvitationGetByCodeOperation {
	return &InvitationGetByCodeOperation{
		repo: r,
		code: code,
	}
}

func (op *InvitationGetByCodeOperation) WithTx(tx *gorm.DB) *InvitationGetByCodeOperation {
	op.tx = tx
	return op
}

func (op *InvitationGetByCodeOperation) Exec() (*model.Invitation, error) {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Where("code = ?", op.code).First(&invitation).Error
	} else {
		err = storage.DB.Where("code = ?", op.code).First(&invitation).Error
	}

	return &invitation, err
}

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

type InvitationListOperation struct {
	repo *InvitationRepository
	tx   *gorm.DB
}

func (r *InvitationRepository) List() *InvitationListOperation {
	return &InvitationListOperation{
		repo: r,
	}
}

func (op *InvitationListOperation) WithTx(tx *gorm.DB) *InvitationListOperation {
	op.tx = tx
	return op
}

func (op *InvitationListOperation) Exec() ([]model.Invitation, error) {
	var invitations []model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Find(&invitations).Error
	} else {
		err = storage.DB.Find(&invitations).Error
	}

	return invitations, err
}

type InvitationListByCreatedByOperation struct {
	repo      *InvitationRepository
	createdBy int64
	tx        *gorm.DB
}

func (r *InvitationRepository) ListByCreatedBy(createdBy int64) *InvitationListByCreatedByOperation {
	return &InvitationListByCreatedByOperation{
		repo:      r,
		createdBy: createdBy,
	}
}

func (op *InvitationListByCreatedByOperation) WithTx(tx *gorm.DB) *InvitationListByCreatedByOperation {
	op.tx = tx
	return op
}

func (op *InvitationListByCreatedByOperation) Exec() ([]model.Invitation, error) {
	var invitations []model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Where("created_by = ?", op.createdBy).Find(&invitations).Error
	} else {
		err = storage.DB.Where("created_by = ?", op.createdBy).Find(&invitations).Error
	}

	return invitations, err
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
	now := time.Now()
	if op.tx != nil {
		return op.tx.Model(&model.Invitation{}).Where("id = ?", op.id).Update("deleted_at", &now).Error
	}
	return storage.DB.Model(&model.Invitation{}).Where("id = ?", op.id).Update("deleted_at", &now).Error
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

type InvitationValidateAndUseOperation struct {
	repo *InvitationRepository
	code string
	tx   *gorm.DB
}

func (r *InvitationRepository) ValidateAndUse(code string) *InvitationValidateAndUseOperation {
	return &InvitationValidateAndUseOperation{
		repo: r,
		code: code,
	}
}

func (op *InvitationValidateAndUseOperation) WithTx(tx *gorm.DB) *InvitationValidateAndUseOperation {
	op.tx = tx
	return op
}

func (op *InvitationValidateAndUseOperation) Exec() error {
	var invitation model.Invitation
	var err error

	if op.tx != nil {
		err = op.tx.Where("code = ? AND is_used = ?", op.code, false).First(&invitation).Error
	} else {
		err = storage.DB.Where("code = ? AND is_used = ?", op.code, false).First(&invitation).Error
	}

	if err != nil {
		return err
	}

	if op.tx != nil {
		return op.tx.Model(&invitation).Updates(map[string]interface{}{
			"is_used": true,
			"used_at": time.Now(),
		}).Error
	}

	return storage.DB.Model(&invitation).Updates(map[string]interface{}{
		"is_used": true,
		"used_at": time.Now(),
	}).Error
}
