package invitation_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
