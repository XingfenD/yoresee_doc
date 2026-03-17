package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserUpdateOperation struct {
	repo *UserRepository
	user *model.User
	tx   *gorm.DB
}

func (r *UserRepository) Update(user *model.User) *UserUpdateOperation {
	return &UserUpdateOperation{
		repo: r,
		user: user,
	}
}

func (op *UserUpdateOperation) WithTx(tx *gorm.DB) *UserUpdateOperation {
	op.tx = tx
	return op
}

func (op *UserUpdateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Save(op.user).Error
	}
	return storage.DB.Save(op.user).Error
}
