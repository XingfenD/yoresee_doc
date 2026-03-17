package user_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
)

type UserService struct {
	userRepo *user_repo.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: user_repo.UserRepo,
	}
}

func (s *UserService) GetByExternalID(externalID string) (*dto.UserResponse, error) {
	userModel, err := s.userRepo.GetByExternalID(externalID).Exec()
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponseFromModel(userModel), nil
}

func (s *UserService) GetByID(id int64) (*dto.UserResponse, error) {
	userModel, err := s.userRepo.GetByID(id).Exec()
	if err != nil {
		return nil, err
	}
	return dto.NewUserResponseFromModel(userModel), nil
}

func (s *UserService) GetIDByExternalID(externalID string) (int64, error) {
	return s.userRepo.GetIDByExternalID(externalID).Exec()
}

var UserSvc = NewUserService()
