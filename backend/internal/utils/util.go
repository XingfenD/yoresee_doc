package utils

import (
	"strings"
)

func Of[T any](v T) *T {
	return &v
}

func GenConfigKey(parts ...string) string {
	return strings.Join(parts, ".")
}
