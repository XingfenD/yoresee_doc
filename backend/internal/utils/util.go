package utils

import (
	"os"
	"strings"
)

func Of[T any](v T) *T {
	return &v
}

func GenConfigKey(parts ...string) string {
	return strings.Join(parts, ".")
}

func GetEnvVar(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return strings.TrimSpace(defaultValue)
	}
	return value
}
