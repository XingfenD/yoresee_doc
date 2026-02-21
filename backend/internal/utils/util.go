package utils

import "github.com/bwmarrin/snowflake"

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
