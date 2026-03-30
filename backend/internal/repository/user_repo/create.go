package user_repo

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserCreateOperation struct {
	repo *UserRepository
	user *model.User
	tx   *gorm.DB
}

func (r *UserRepository) Create(user *model.User) *UserCreateOperation {
	return &UserCreateOperation{
		repo: r,
		user: user,
	}
}

func (op *UserCreateOperation) WithTx(tx *gorm.DB) *UserCreateOperation {
	op.tx = tx
	return op
}

func (op *UserCreateOperation) Exec() error {
	if op.tx != nil {
		if err := op.tx.Create(op.user).Error; err != nil {
			return err
		}
		return op.clearQueryCache()
	}
	if err := storage.DB.Create(op.user).Error; err != nil {
		return err
	}
	return op.clearQueryCache()
}

func (op *UserCreateOperation) clearQueryCache() error {
	_, _ = storage.KVS.Incr(context.Background(), key.KeyUserQueryVersion()).Result()
	return nil
}
