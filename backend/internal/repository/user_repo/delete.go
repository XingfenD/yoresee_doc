package user_repo

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
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
		if err := op.tx.Delete(&model.User{}, op.id).Error; err != nil {
			return err
		}
		return op.clearQueryCache()
	}
	if err := op.repo.db.Delete(&model.User{}, op.id).Error; err != nil {
		return err
	}
	return op.clearQueryCache()
}

func (op *UserDeleteOperation) clearQueryCache() error {
	_, _ = op.repo.redis.Incr(context.Background(), key.KeyUserQueryVersion()).Result()
	return nil
}
