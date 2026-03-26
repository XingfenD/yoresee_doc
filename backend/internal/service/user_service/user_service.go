package user_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
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
		logrus.Errorf("[Service layer: UserService] GetByExternalID failed, external_id=%s, err=%+v", externalID, err)
		return nil, status.GenErrWithCustomMsg(status.StatusUserNotFound, "user not found")
	}
	return dto.NewUserResponseFromModel(userModel), nil
}

func (s *UserService) GetByID(id int64) (*dto.UserResponse, error) {
	userModel, err := s.userRepo.GetByID(id).Exec()
	if err != nil {
		logrus.Errorf("[Service layer: UserService] GetByID failed, id=%d, err=%+v", id, err)
		return nil, status.GenErrWithCustomMsg(status.StatusUserNotFound, "user not found")
	}
	return dto.NewUserResponseFromModel(userModel), nil
}

func (s *UserService) GetIDByExternalID(externalID string) (int64, error) {
	id, err := s.userRepo.GetIDByExternalID(externalID).Exec()
	if err != nil {
		logrus.Errorf("[Service layer: UserService] GetIDByExternalID failed, external_id=%s, err=%+v", externalID, err)
		return 0, status.GenErrWithCustomMsg(status.StatusUserNotFound, "user not found")
	}
	return id, nil
}

var UserSvc = NewUserService()
