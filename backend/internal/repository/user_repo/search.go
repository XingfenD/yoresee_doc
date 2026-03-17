package user_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserSearchOperation struct {
	repo  *UserRepository
	query string
	tx    *gorm.DB
}

func (r *UserRepository) Search(query string) *UserSearchOperation {
	return &UserSearchOperation{
		repo:  r,
		query: query,
	}
}

func (op *UserSearchOperation) WithTx(tx *gorm.DB) *UserSearchOperation {
	op.tx = tx
	return op
}

func (op *UserSearchOperation) Exec() ([]model.User, error) {
	var users []model.User
	var err error

	if op.tx != nil {
		db := op.tx.Where("username LIKE ? OR email LIKE ?", "%"+op.query+"%", "%"+op.query+"%")
		err = db.Or("CAST(id AS CHAR) = ?", op.query).Find(&users).Error
	} else {
		db := storage.DB.Where("username LIKE ? OR email LIKE ?", "%"+op.query+"%", "%"+op.query+"%")
		err = db.Or("CAST(id AS CHAR) = ?", op.query).Find(&users).Error
	}

	return users, err
}
