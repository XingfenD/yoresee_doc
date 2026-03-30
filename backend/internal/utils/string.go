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
