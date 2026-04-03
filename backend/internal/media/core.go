package media

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func normalizeMIMEType(contentType string) string {
	normalized := strings.ToLower(strings.TrimSpace(contentType))
	if idx := strings.Index(normalized, ";"); idx > 0 {
		normalized = strings.TrimSpace(normalized[:idx])
	}
	if normalized == "image/jpg" {
		return "image/jpeg"
	}
	return normalized
}

func NormalizeContentTypeByAllowed(file []byte, provided string, allowed map[string]string) string {
	normalized := normalizeMIMEType(provided)
	if _, ok := allowed[normalized]; ok {
		return normalized
	}

	detected := normalizeMIMEType(http.DetectContentType(file))
	if _, ok := allowed[detected]; ok {
		return detected
	}
	return ""
}

func IsContentTypeAllowed(contentType string, allowed map[string]string) bool {
	_, ok := allowed[normalizeMIMEType(contentType)]
	return ok
}

func ResolveExtension(filename, contentType string, allowed map[string]string, defaultExt string) string {
	ext := strings.TrimSpace(strings.ToLower(filepath.Ext(filename)))
	if ext != "" {
		for _, candidate := range allowed {
			if ext == candidate {
				return ext
			}
		}
	}

	contentType = normalizeMIMEType(contentType)
	if mapped, ok := allowed[contentType]; ok {
		return mapped
	}

	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		for _, candidate := range exts {
			candidate = strings.TrimSpace(strings.ToLower(candidate))
			for _, allowedExt := range allowed {
				if candidate == allowedExt {
					return candidate
				}
			}
		}
	}
	return strings.TrimSpace(strings.ToLower(defaultExt))
}

func BuildVersionedObjectKey(prefix, ownerExternalID string, version int64, ext string) string {
	prefix = strings.Trim(strings.TrimSpace(prefix), "/")
	if prefix == "" {
		return fmt.Sprintf("%s/v%d%s", ownerExternalID, version, ext)
	}
	return fmt.Sprintf("%s/%s/v%d%s", prefix, ownerExternalID, version, ext)
}

func BuildStorageURL(objectKey string) string {
	if strings.TrimSpace(objectKey) == "" {
		return ""
	}
	return storage.BuildPublicObjectPath(objectKey)
}
