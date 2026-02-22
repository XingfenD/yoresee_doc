package service

import (
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.UserRepo,
	}
}

func (s *AuthService) Register(userCreate *dto.UserCreate) error {
	return utils.WithTransaction(func(tx *gorm.DB) error {
		mode := ConfigSvc.GetSystemRegisterMode()
		if mode == constant.RegisterMode_Invite {
			if userCreate.InvitationCode == nil || *userCreate.InvitationCode == "" {
				return status.GenErrWithCustomMsg(status.StatusParamError, "the system enable invitation mode, but the invitation code is empty")
			}
			err := InvitationSvc.ValidateAndUseWithTx(tx, *userCreate.InvitationCode)
			if err != nil {
				return err
			}
		}

		_, err := s.userRepo.GetByEmail(userCreate.Email).WithTx(tx).Exec()
		if err == nil {
			return status.GenErrWithCustomMsg(status.StatusUserAlreadyExists, "email already registered")
		}

		hashedPwd, err := utils.HashPassword(userCreate.Password)
		if err != nil {
			return err
		}

		user := &model.User{
			ExternalID:     utils.GenerateExternalID(utils.ExternalIDContextUser),
			Username:       userCreate.Username,
			PasswordHash:   hashedPwd,
			Email:          userCreate.Email,
			Status:         1,
			InvitationCode: userCreate.InvitationCode,
		}

		if err := s.userRepo.Create(user).WithTx(tx).Exec(); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) Login(email string, password string) (string, *dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(email).Exec()
	if err != nil {
		return "", nil, status.StatusUserNotFound
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", nil, status.StatusInvalidPassword
	}

	token, err := utils.GenerateToken(user.ExternalID, user.Username)
	if err != nil {
		return "", nil, err
	}

	userResponse := dto.NewUserResponseFromModel(user)

	return token, userResponse, nil
}

var AuthSvc = NewAuthService()
