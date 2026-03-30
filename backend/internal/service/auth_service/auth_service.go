package auth_service

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/auth"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/user_repo"
	"github.com/XingfenD/yoresee_doc/internal/service/config_service"
	"github.com/XingfenD/yoresee_doc/internal/service/invitation_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *user_repo.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: user_repo.UserRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, userCreate *dto.UserCreate) error {
	return utils.WithTransaction(func(tx *gorm.DB) error {
		mode := config_service.ConfigSvc.GetSystemRegisterMode(ctx)
		if mode == constant.RegisterMode_Invite {
			if userCreate.InvitationCode == nil || *userCreate.InvitationCode == "" {
				return status.GenErrWithCustomMsg(status.StatusParamError, "the system enable invitation mode, but the invitation code is empty")
			}
			err := invitation_service.InvitationSvc.ValidateAndUseWithTx(tx, *userCreate.InvitationCode)
			if err != nil {
				logrus.Errorf("[Service layer: AuthService] ValidateAndUseWithTx failed, invitation_code=%s, err=%+v", *userCreate.InvitationCode, err)
				_ = invitation_service.InvitationSvc.CreateInvitationRecord(&model.InvitationRecord{
					Code:   *userCreate.InvitationCode,
					UsedBy: userCreate.Email,
					Status: "failed",
					UsedAt: time.Now(),
				})
				return status.GenErrWithCustomMsg(err, "invitation validation failed")
			}
		}

		_, err := s.userRepo.GetByEmail(userCreate.Email).WithTx(tx).Exec()
		if err == nil {
			return status.GenErrWithCustomMsg(status.StatusUserAlreadyExists, "email already registered")
		}

		hashedPwd, err := utils.HashPassword(userCreate.Password)
		if err != nil {
			logrus.Errorf("[Service layer: AuthService] HashPassword failed, email=%s, err=%+v", userCreate.Email, err)
			return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "hash password failed")
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
			logrus.Errorf("[Service layer: AuthService] Create user failed, email=%s, err=%+v", userCreate.Email, err)
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create user failed")
		}

		if mode == constant.RegisterMode_Invite && userCreate.InvitationCode != nil && *userCreate.InvitationCode != "" {
			_ = invitation_service.InvitationSvc.CreateInvitationRecordWithTx(tx, &model.InvitationRecord{
				Code:         *userCreate.InvitationCode,
				UsedByUserID: utils.Of(user.ID),
				UsedBy:       user.Email,
				Status:       "success",
				UsedAt:       time.Now(),
			})
		}

		return nil
	})
}

func (s *AuthService) Login(email string, password string) (string, *dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(email).Exec()
	if err != nil {
		return "", nil, status.StatusUserNotFound
	}
	if user.Status <= 0 {
		return "", nil, status.GenErrWithCustomMsg(status.StatusPermissionDenied, "user is banned")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", nil, status.StatusInvalidPassword
	}

	token, err := auth.GenerateToken(user.ExternalID, user.Username)
	if err != nil {
		logrus.Errorf("[Service layer: AuthService] GenerateToken failed, user_external_id=%s, err=%+v", user.ExternalID, err)
		return "", nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "generate token failed")
	}

	userResponse := dto.NewUserResponseFromModel(user)

	return token, userResponse, nil
}

var AuthSvc = NewAuthService()
