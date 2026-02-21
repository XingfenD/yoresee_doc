package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type UserRepository struct{}

var UserRepo = &UserRepository{}

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
		return op.tx.Create(op.user).Error
	}
	return storage.DB.Create(op.user).Error
}

type UserGetByUsernameOperation struct {
	repo     *UserRepository
	username string
	tx       *gorm.DB
}

func (r *UserRepository) GetByUsername(username string) *UserGetByUsernameOperation {
	return &UserGetByUsernameOperation{
		repo:     r,
		username: username,
	}
}

func (op *UserGetByUsernameOperation) WithTx(tx *gorm.DB) *UserGetByUsernameOperation {
	op.tx = tx
	return op
}

func (op *UserGetByUsernameOperation) Exec() (*model.User, error) {
	var user model.User
	var err error

	if op.tx != nil {
		err = op.tx.Where("username = ?", op.username).First(&user).Error
	} else {
		err = storage.DB.Where("username = ?", op.username).First(&user).Error
	}

	return &user, err
}

type UserGetByIDOperation struct {
	repo *UserRepository
	id   int64
	tx   *gorm.DB
}

func (r *UserRepository) GetByID(id int64) *UserGetByIDOperation {
	return &UserGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *UserGetByIDOperation) WithTx(tx *gorm.DB) *UserGetByIDOperation {
	op.tx = tx
	return op
}

func (op *UserGetByIDOperation) Exec() (*model.User, error) {
	var user model.User
	var err error

	if op.tx != nil {
		err = op.tx.First(&user, op.id).Error
	} else {
		err = storage.DB.First(&user, op.id).Error
	}

	return &user, err
}

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
