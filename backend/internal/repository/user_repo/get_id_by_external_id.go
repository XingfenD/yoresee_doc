package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserGetIDByExternalIDOperation struct {
	repo       *UserRepository
	externalID string
	tx         *gorm.DB
}

func (r *UserRepository) GetIDByExternalID(externalID string) *UserGetIDByExternalIDOperation {
	return &UserGetIDByExternalIDOperation{
		repo:       r,
		externalID: externalID,
	}
}

func (op *UserGetIDByExternalIDOperation) WithTx(tx *gorm.DB) *UserGetIDByExternalIDOperation {
	op.tx = tx
	return op
}

func (op *UserGetIDByExternalIDOperation) Exec() (int64, error) {
	var id int64
	var err error

	if op.tx != nil {
		err = op.tx.Model(&model.User{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	} else {
		err = storage.DB.Model(&model.User{}).Where("external_id = ?", op.externalID).Pluck("id", &id).Error
	}

	return id, err
}
