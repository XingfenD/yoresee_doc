package utils

import (
	"crypto/md5"
	"encoding/hex"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node     *snowflake.Node
	nodeOnce sync.Once
	nodeErr  error
)

type ExternalIDContext string

const (
	ExternalIDContextUser         ExternalIDContext = "user"
	ExternalIDKnowledgeBase       ExternalIDContext = "knowledge_base"
	ExternalIDContextDocument     ExternalIDContext = "document"
	ExternalIDContextUserGroup    ExternalIDContext = "user_group"
	ExternalIDContextOrgNode      ExternalIDContext = "org_node"
	ExternalIDContextComment      ExternalIDContext = "comment"
	ExternalIDContextNotification ExternalIDContext = "notification"
)

func GenerateExternalID(context ExternalIDContext) string {
	initSnowflakeNode()
	if nodeErr != nil {
		panic(nodeErr)
	}
	id := node.Generate()
	combined := string(context) + ":" + id.String()
	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func initSnowflakeNode() {
	nodeOnce.Do(func() {
		node, nodeErr = snowflake.NewNode(1)
	})
}
