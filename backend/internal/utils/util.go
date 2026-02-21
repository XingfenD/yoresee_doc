package utils

import (
	"strings"

	"github.com/bwmarrin/snowflake"
)

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func Of[T any](v T) *T {
	return &v
}

func GenConfigKey(parts ...string) string {
	return strings.Join(parts, ".")
}
