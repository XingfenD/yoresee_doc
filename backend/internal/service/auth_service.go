package service

import (
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

func (s *AuthService) Register(username string, password string, email string, invitationCode *string) error {
	return utils.WithTransaction(func(tx *gorm.DB) error {
		mode, _ := ConfigSvc.Get("registration_mode")
		if mode == "invitation" {
			if invitationCode == nil || *invitationCode == "" {
				return status.GenErrWithCustomMsg(status.StatusParamError, "the system enable invitation mode, but the invitation code is empty")
			}
			err := InvitationSvc.ValidateAndUseWithTx(tx, *invitationCode)
			if err != nil {
				return err
			}
		}

		// 使用新的链式调用API
		_, err := s.userRepo.GetByUsername(username).WithTx(tx).Exec()
		if err == nil {
			return status.StatusUserAlreadyExists
		}

		hashedPwd, err := utils.HashPassword(password)
		if err != nil {
			return err
		}

		user := &model.User{
			Username:       username,
			PasswordHash:   hashedPwd,
			Email:          email,
			RoleID:         2, // default to normal user (assuming 2 is user role)
			Status:         1,
			InvitationCode: invitationCode,
		}

		// 使用新的链式调用API
		if err := s.userRepo.Create(user).WithTx(tx).Exec(); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) Login(username string, password string) error {
	user, err := s.userRepo.GetByUsername(username).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return status.StatusInvalidPassword
	}

	return nil
}

var AuthSvc = NewAuthService()
