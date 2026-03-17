package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type ListUserByExternalOperation struct {
	repo           *UserRepository
	externalIDList []string
	tx             *gorm.DB
}

func (r *UserRepository) ListByExternal(externalIDList []string) *ListUserByExternalOperation {
	return &ListUserByExternalOperation{
		repo:           r,
		externalIDList: externalIDList,
	}
}

func (op *ListUserByExternalOperation) WithTx(tx *gorm.DB) *ListUserByExternalOperation {
	op.tx = tx
	return op
}

func (op *ListUserByExternalOperation) Exec() ([]model.User, error) {
	var users []model.User
	var err error

	if op.tx != nil {
		err = op.tx.Where("external_id IN ?", op.externalIDList).Find(&users).Error
	} else {
		err = storage.DB.Where("external_id IN ?", op.externalIDList).Find(&users).Error
	}

	return users, err
}
