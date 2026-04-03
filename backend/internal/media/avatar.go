package media

import (
	"strings"
)

const (
	MaxAvatarSize          = 5 * 1024 * 1024
	defaultAvatarExtension = ".jpg"
)

var avatarContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

func NormalizeAvatarContentType(file []byte, provided string) string {
	return NormalizeContentTypeByAllowed(file, provided, avatarContentTypes)
}

func IsSupportedAvatarContentType(contentType string) bool {
	return IsContentTypeAllowed(contentType, avatarContentTypes)
}

func ResolveAvatarExt(filename, contentType string) string {
	return ResolveExtension(filename, contentType, avatarContentTypes, defaultAvatarExtension)
}

func BuildAvatarObjectKey(userExternalID string, avatarVersion int64, ext string) string {
	return BuildVersionedObjectKey("avatars", userExternalID, avatarVersion, ext)
}

func BuildAvatarURL(userExternalID, objectKey string, avatarVersion int64) string {
	if strings.TrimSpace(userExternalID) == "" || strings.TrimSpace(objectKey) == "" {
		return ""
	}
	return BuildStorageURL(objectKey)
}
