package invitation_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

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
