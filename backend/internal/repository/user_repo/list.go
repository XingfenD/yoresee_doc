package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserListOperation struct {
	repo *UserRepository
	tx   *gorm.DB
}

func (r *UserRepository) List() *UserListOperation {
	return &UserListOperation{
		repo: r,
	}
}

func (op *UserListOperation) WithTx(tx *gorm.DB) *UserListOperation {
	op.tx = tx
	return op
}

func (op *UserListOperation) Exec() ([]model.User, error) {
	var users []model.User
	var err error

	if op.tx != nil {
		err = op.tx.Find(&users).Error
	} else {
		err = storage.DB.Find(&users).Error
	}

	return users, err
}
