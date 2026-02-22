package model

import (
	"time"
)

// Permission 权限枚举
type Permission string

const (
	// 文档类权限
	PermissionRead   Permission = "read"    // 可阅读
	PermissionNoRead Permission = "no_read" // 不可阅读（显式拒绝）
	PermissionEdit   Permission = "edit"    // 可编辑（增删改内容）
	PermissionManage Permission = "manage"  // 可管理（元数据、权限、删除）

	// 组织类权限
	PermissionAdmin       Permission = "admin"        // 完全管理（增删改用户/组织）
	PermissionEditMembers Permission = "edit_members" // 成员管理（添加/移除成员）
	PermissionViewMembers Permission = "view_members" // 查看成员列表

	// 系统权限
	PermissionCreate   Permission = "create"   // 创建子资源
	PermissionTransfer Permission = "transfer" // 转移所有权
	PermissionAudit    Permission = "audit"    // 审计查看
)

// ScopeType 权限范围类型
type ScopeType string

const (
	ScopeTypeExact     ScopeType = "exact"     // 仅当前资源
	ScopeTypeChildren  ScopeType = "children"  // 直接子资源
	ScopeTypeTree      ScopeType = "tree"      // 当前+所有嵌套子文档（用于DOCUMENT_TREE）
	ScopeTypeRecursive ScopeType = "recursive" // 递归所有层级（用于组织架构）
)

// PermissionRule 权限规则模型
type PermissionRule struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`

	// 资源维度
	ResourceType ResourceType `gorm:"size:32;not null;index:idx_resource,priority:2" json:"resource_type"`
	ResourceID   string       `gorm:"size:64;not null;index:idx_resource,priority:3" json:"resource_id"`
	ResourcePath string       `gorm:"type:ltree" json:"resource_path"` // 用于树形权限范围

	// 主体维度
	SubjectType SubjectType `gorm:"size:32;not null;index:idx_subject,priority:2" json:"subject_type"`
	SubjectID   string      `gorm:"size:64;not null;index:idx_subject,priority:3" json:"subject_id"`

	// 权限内容
	Permissions []string  `gorm:"type:varchar(32)[];not null" json:"permissions"` // 权限列表
	ScopeType   ScopeType `gorm:"size:32;default:exact" json:"scope_type"`        // 权限范围类型
	IsDeny      bool      `gorm:"default:false" json:"is_deny"`                   // 是否显式拒绝
	Priority    int       `gorm:"default:100" json:"priority"`                    // 优先级（数字越小优先级越高）

	// 有效期
	ValidFrom  *time.Time `json:"valid_from"`  // 生效时间
	ValidUntil *time.Time `json:"valid_until"` // 失效时间

	// 审计信息
	CreatedBy string    `gorm:"size:64" json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (PermissionRule) TableName() string {
	return "permission_rules"
}
