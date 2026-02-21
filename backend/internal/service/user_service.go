package service

import (
	"github.com/XingfenD/yoresee_doc/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.UserRepo,
	}
}

var UserSvc = NewUserService()
