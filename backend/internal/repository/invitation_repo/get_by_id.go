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
