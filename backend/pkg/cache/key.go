package cache

import (
	"fmt"
)

func KeyDocTree(userID, kbID *int64) string {
	if kbID != nil {
		return fmt.Sprintf("doc:tree:kb:%d", *kbID)
	}
	return fmt.Sprintf("doc:tree:u:%d", userID)
}

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

func KeySystemConfig(configKey string) string {
	return fmt.Sprintf("dms:config:%s", configKey)
}
