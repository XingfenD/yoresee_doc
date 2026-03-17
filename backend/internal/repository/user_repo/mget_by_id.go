package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type MGetUserByIDOperation struct {
	repo    *UserRepository
	userIDs []int64
	tx      *gorm.DB
}

func (r *UserRepository) MGetUserByID(userIDs []int64) *MGetUserByIDOperation {
	return &MGetUserByIDOperation{
		repo:    r,
		userIDs: userIDs,
	}
}

func (op *MGetUserByIDOperation) WithTx(tx *gorm.DB) *MGetUserByIDOperation {
	op.tx = tx
	return op
}

func (op *MGetUserByIDOperation) Exec() (map[int64]*model.User, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}
	result := make(map[int64]*model.User)

	if len(op.userIDs) == 0 {
		return result, nil
	}

	var users []model.User
	var err error

	err = op.tx.Where("id IN ?", op.userIDs).Find(&users).Error

	if err != nil {
		return nil, err
	}

	for _, user := range users {
		userCopy := user
		result[user.ID] = &userCopy
	}

	return result, nil
}
