package auth_service

import (
	"bytes"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

const maxAvatarSize = 5 * 1024 * 1024

func (s *AuthService) UpdateProfile(userExternalID string, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	if strings.TrimSpace(userExternalID) == "" || req == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user_external_id and request are required")
	}
	if req.Username == nil && req.Email == nil && req.Nickname == nil && req.Password == nil && req.Avatar == nil && req.AvatarFile == nil {
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

	if req.AvatarFile != nil {
		avatarURL, uploadErr := uploadAvatar(user.ExternalID, req.AvatarFile, req.AvatarFilename, req.AvatarContentType)
		if uploadErr != nil {
			return nil, status.GenErrWithCustomMsg(status.StatusServiceInternalError, "upload avatar failed")
		}
		user.Avatar = avatarURL
	} else if req.Avatar != nil {
		user.Avatar = strings.TrimSpace(*req.Avatar)
	}

	if err := s.userRepo.Update(user).Exec(); err != nil {
		return nil, status.StatusWriteDBError
	}

	return dto.NewUserResponseFromModel(user), nil
}

func uploadAvatar(userExternalID string, avatarFile []byte, filename, contentType *string) (string, error) {
	if len(avatarFile) == 0 {
		return "", status.GenErrWithCustomMsg(status.StatusParamError, "avatar file is empty")
	}
	if len(avatarFile) > maxAvatarSize {
		return "", status.GenErrWithCustomMsg(status.StatusParamError, "avatar file exceeds 5MB limit")
	}
	if storage.MinioClient == nil {
		return "", status.GenErrWithCustomMsg(status.StatusServiceInternalError, "minio client is not initialized")
	}

	fileName := "avatar"
	if filename != nil {
		fileName = strings.TrimSpace(*filename)
	}

	fileContentType := ""
	if contentType != nil {
		fileContentType = strings.TrimSpace(*contentType)
	}
	if fileContentType == "" {
		fileContentType = http.DetectContentType(avatarFile)
	}

	objectName := fmt.Sprintf(
		"avatars/%s/%d%s",
		userExternalID,
		time.Now().UnixNano(),
		resolveAvatarExt(fileName, fileContentType),
	)

	return storage.UploadFile(
		config.GlobalConfig.Minio.Bucket,
		objectName,
		bytes.NewReader(avatarFile),
		int64(len(avatarFile)),
		fileContentType,
	)
}

func resolveAvatarExt(filename, contentType string) string {
	ext := strings.TrimSpace(strings.ToLower(filepath.Ext(filename)))
	if ext != "" {
		return ext
	}
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}
