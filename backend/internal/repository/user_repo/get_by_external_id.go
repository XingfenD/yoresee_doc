package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserGetByExternalIDOperation struct {
	repo       *UserRepository
	externalID string
	tx         *gorm.DB
}

func (r *UserRepository) GetByExternalID(externalID string) *UserGetByExternalIDOperation {
	return &UserGetByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *UserGetByExternalIDOperation) WithTx(tx *gorm.DB) *UserGetByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *UserGetByExternalIDOperation) Exec() (*model.User, error) {
	var user model.User
	var err error

	if op.tx != nil {
		err = op.tx.Where("external_id = ?", op.externalID).First(&user).Error
	} else {
		err = storage.DB.Where("external_id = ?", op.externalID).First(&user).Error
	}

	return &user, err
}
