package model

type RoleScope int64

const (
	RoleScope_AllUser   RoleScope = 0 // 全局角色，默认适用于所有用户
	RoleScope_OrgNode   RoleScope = 1 // 组织角色，仅适用于特定组织内的用户
	RoleScope_User      RoleScope = 2 // 用户角色，仅适用于特定用户
	RoleScope_UserGroup RoleScope = 3 // 用户组角色，仅适用于特定用户组
)

type Role struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"size:64;not null;unique" json:"name"`
	ParentRoleID *int64 `json:"parent_role_id"`
	IsSystem     bool   `gorm:"not null;default:false" json:"is_system"`
	CreatedBy    *int64 `json:"created_by"`
	CreatedAt    *int64 `json:"created_at"`
}

func (Role) TableName() string {
	return "roles"
}
