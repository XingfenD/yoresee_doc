package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"gorm.io/gorm"
)

type UserGetByEmailOperation struct {
	repo  *UserRepository
	email string
	tx    *gorm.DB
}

func (r *UserRepository) GetByEmail(email string) *UserGetByEmailOperation {
	return &UserGetByEmailOperation{
		repo:  r,
		email: email,
	}
}

func (op *UserGetByEmailOperation) WithTx(tx *gorm.DB) *UserGetByEmailOperation {
	op.tx = tx
	return op
}

func (op *UserGetByEmailOperation) Exec() (*model.User, error) {
	var user model.User
	var err error

	if op.tx != nil {
		err = op.tx.Where("email = ?", op.email).First(&user).Error
	} else {
		err = op.repo.db.Where("email = ?", op.email).First(&user).Error
	}

	return &user, err
}
