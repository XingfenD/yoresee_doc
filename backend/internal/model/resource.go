package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ResourceType string

const (
	ResourceTypeKnowledgeBase ResourceType = "knowledge_base" // 知识库（根级容器）
	ResourceTypeDocument      ResourceType = "document"       // 单个文档
	ResourceTypeDocumentTree  ResourceType = "document_tree"  // 文档及其嵌套子文档树

	ResourceTypeOrgStructure ResourceType = "org_structure" // 组织架构节点
	ResourceTypeUserGroup    ResourceType = "user_group"    // 用户组
	ResourceTypeUser         ResourceType = "user"          // 单一用户
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type Resource struct {
	ID       string       `gorm:"primaryKey;size:64" json:"id"`
	Type     ResourceType `gorm:"size:32;not null" json:"type"`
	ParentID *string      `gorm:"size:64" json:"parent_id"`
	Path     string       `gorm:"type:ltree" json:"path"` // PostgreSQL层级路径类型
	Metadata JSONB        `gorm:"type:jsonb" json:"metadata"`
}

func (Resource) TableName() string {
	return "resources"
}
