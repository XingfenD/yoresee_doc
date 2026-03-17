package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserDeleteOperation struct {
	repo *UserRepository
	id   int64
	tx   *gorm.DB
}

func (r *UserRepository) Delete(id int64) *UserDeleteOperation {
	return &UserDeleteOperation{
		repo: r,
		id:   id,
	}
}

func (op *UserDeleteOperation) WithTx(tx *gorm.DB) *UserDeleteOperation {
	op.tx = tx
	return op
}

func (op *UserDeleteOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Delete(&model.User{}, op.id).Error
	}
	return storage.DB.Delete(&model.User{}, op.id).Error
}
