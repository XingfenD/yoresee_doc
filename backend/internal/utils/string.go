package utils

import (
	"strconv"
	"strings"
)

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ParseInt64(value string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
}

func NormalizeToken(raw, fallback string) string {
	fallback = strings.TrimSpace(fallback)
	if fallback == "" {
		fallback = "default"
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}

	var b strings.Builder
	for _, r := range raw {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '-' || r == '_' {
			b.WriteRune(r)
			continue
		}
		b.WriteByte('-')
	}

	if b.Len() == 0 {
		return fallback
	}
	return b.String()
}
