package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

type ExternalIDContext string

const (
	ExternalIDContextUser      ExternalIDContext = "user"
	ExternalIDContextDocument  ExternalIDContext = "document"
	ExternalIDContextUserGroup ExternalIDContext = "user_group"
	ExternalIDContextOrgNode   ExternalIDContext = "org_node"
)

func GenerateExternalID(context ExternalIDContext) string {
	id := node.Generate()
	combined := string(context) + ":" + id.String()
	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}
