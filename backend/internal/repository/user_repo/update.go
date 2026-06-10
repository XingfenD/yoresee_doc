package user_repo

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
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
		if err := op.tx.Save(op.user).Error; err != nil {
			return err
		}
		return op.clearQueryCache()
	}
	if err := op.repo.db.Save(op.user).Error; err != nil {
		return err
	}
	return op.clearQueryCache()
}

func (op *UserUpdateOperation) clearQueryCache() error {
	_, _ = op.repo.redis.Incr(context.Background(), key.KeyUserQueryVersion()).Result()
	return nil
}
