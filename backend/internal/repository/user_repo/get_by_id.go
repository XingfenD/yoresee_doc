package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserGetByIDOperation struct {
	repo *UserRepository
	id   int64
	tx   *gorm.DB
}

func (r *UserRepository) GetByID(id int64) *UserGetByIDOperation {
	return &UserGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *UserGetByIDOperation) WithTx(tx *gorm.DB) *UserGetByIDOperation {
	op.tx = tx
	return op
}

func (op *UserGetByIDOperation) Exec() (*model.User, error) {
	var user model.User
	var err error

	if op.tx != nil {
		err = op.tx.First(&user, op.id).Error
	} else {
		err = storage.DB.First(&user, op.id).Error
	}

	return &user, err
}
