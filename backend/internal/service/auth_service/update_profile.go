package auth_service

import (
	"bytes"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/media"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func (s *AuthService) UpdateProfile(userExternalID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	if strings.TrimSpace(userExternalID) == "" || req == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user_external_id and request are required")
	}
	if req.Username == nil && req.Email == nil && req.Nickname == nil && req.Password == nil && req.AvatarFile == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "no profile fields to update")
	}

	user, err := s.userRepo.GetByExternalID(userExternalID).Exec()
	if err != nil {
		return nil, status.GenErrWithCustomMsg(status.StatusUserNotFound, "user not found")
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if username == "" {
			return nil, status.GenErrWithCustomMsg(status.StatusParamError, "username cannot be empty")
		}
		user.Username = username
	}

	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" {
			return nil, status.GenErrWithCustomMsg(status.StatusParamError, "email cannot be empty")
		}
		existingUser, getErr := s.userRepo.GetByEmail(email).Exec()
		if getErr == nil && existingUser.ID != user.ID {
			return nil, status.GenErrWithCustomMsg(status.StatusUserAlreadyExists, "email already registered")
		}
		user.Email = email
	}

	if req.Nickname != nil {
		user.Nickname = strings.TrimSpace(*req.Nickname)
	}

	if req.Password != nil {
		password := strings.TrimSpace(*req.Password)
		if password == "" {
			return nil, status.GenErrWithCustomMsg(status.StatusParamError, "password cannot be empty")
		}
		hashedPwd, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "hash password failed")
		}
		user.PasswordHash = hashedPwd
	}

	var oldAvatarObjectKey string
	if req.AvatarFile != nil {
		oldAvatarObjectKey = user.AvatarObjectKey
		avatarObjectKey, nextVersion, uploadErr := uploadAvatar(
			user.ExternalID,
			user.AvatarVersion+1,
			req.AvatarFile,
			req.AvatarFilename,
			req.AvatarContentType,
		)
		if uploadErr != nil {
			return nil, status.GenErrWithCustomMsg(uploadErr, "upload avatar failed")
		}
		now := time.Now()
		user.AvatarObjectKey = avatarObjectKey
		user.AvatarVersion = nextVersion
		user.AvatarUpdatedAt = &now
		user.Avatar = ""
	}

	if err := s.userRepo.Update(user).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}

	if oldAvatarObjectKey != "" && oldAvatarObjectKey != user.AvatarObjectKey {
		_ = storage.DeleteFile(config.GlobalConfig.Minio.Bucket, oldAvatarObjectKey)
	}

	resp := dto.NewUserResponseFromModel(user)
	resp.Avatar = media.BuildAvatarURL(resp.ExternalID, resp.AvatarObjectKey, resp.AvatarVersion)
	return resp, nil
}

func uploadAvatar(userExternalID string, avatarVersion int64, avatarFile []byte, filename, contentType *string) (string, int64, error) {
	if len(avatarFile) == 0 {
		return "", 0, status.GenErrWithCustomMsg(status.StatusParamError, "avatar file is empty")
	}
	if len(avatarFile) > media.MaxAvatarSize {
		return "", 0, status.GenErrWithCustomMsg(status.StatusParamError, "avatar file exceeds 5MB limit")
	}
	if storage.MinioClient == nil {
		return "", 0, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "minio client is not initialized")
	}

	fileName := "avatar"
	if filename != nil {
		fileName = strings.TrimSpace(*filename)
	}

	fileContentType := "application/octet-stream"
	if contentType != nil {
		fileContentType = strings.TrimSpace(*contentType)
	}
	fileContentType = media.NormalizeAvatarContentType(avatarFile, fileContentType)
	if !media.IsSupportedAvatarContentType(fileContentType) {
		return "", 0, status.GenErrWithCustomMsg(status.StatusParamError, "avatar content type is not supported")
	}

	objectName := media.BuildAvatarObjectKey(userExternalID, avatarVersion, media.ResolveAvatarExt(fileName, fileContentType))

	err := storage.PutFile(
		config.GlobalConfig.Minio.Bucket,
		objectName,
		bytes.NewReader(avatarFile),
		int64(len(avatarFile)),
		fileContentType,
	)
	if err != nil {
		return "", 0, err
	}
	return objectName, avatarVersion, nil
}
