package invitation_repo

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
