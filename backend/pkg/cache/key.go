package cache

import (
	"fmt"
)

type KeyObjectTypeEnum int

const (
	KeyObjectTypeEnum_User KeyObjectTypeEnum = iota
	KeyObjectTypeEnum_Doc
	KeyObjectTypeEnum_KnowledgeBase
)

var objStringMap = map[KeyObjectTypeEnum]string{
	KeyObjectTypeEnum_User:          "user",
	KeyObjectTypeEnum_Doc:           "doc",
	KeyObjectTypeEnum_KnowledgeBase: "doc",
}

func KeyIDByExternalID(obj KeyObjectTypeEnum, externalID string) string {
	for k, v := range objStringMap {
		if k == obj {
			return fmt.Sprintf("dms:%s:extid:%s", v, externalID)
		}
	}

	return fmt.Sprintf("dms:unknownobj_%d:extid:%s", int(obj), externalID)
}

func KeyModelByExternalID(obj KeyObjectTypeEnum, externalID string) string {
	for k, v := range objStringMap {
		if k == obj {
			return fmt.Sprintf("dms:%s:model:%s", v, externalID)
		}
	}

	return fmt.Sprintf("dms:unknownobj_%d:model:%s", int(obj), externalID)
}

func KeyDocSubtreeVersion(path string) string {
	return fmt.Sprintf("dms:doc:version:%s", path)
}

func KeyDocSubtree(path string, version int64, depth *int) string {
	depthKey := "all"
	if depth != nil {
		depthKey = fmt.Sprintf("%d", *depth)
	}
	return fmt.Sprintf("dms:doc:subtree:%s:v%d:%s", path, version, depthKey)
}

func KeyUserQueryList(queryHash string, page, pageSize int) string {
	return fmt.Sprintf("dms:user:query:%s:p%d:s%d", queryHash, page, pageSize)
}

func KeyUserQueryPrefix() string {
	return "dms:user:query:*"
}

func KeyUserQueryVersion() string {
	return "dms:user:query:version"
}
