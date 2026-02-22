package model

// SubjectType 主体类型枚举
type SubjectType string

const (
	SubjectTypeUser      SubjectType = "user"       // 单一用户
	SubjectTypeUserGroup SubjectType = "user_group" // 用户组（扁平）
	SubjectTypeOrgNode   SubjectType = "org_node"   // 组织架构节点
)

// Subject 主体模型
type Subject struct {
	ID         string      `gorm:"primaryKey;size:64" json:"id"`
	Namespace  string      `gorm:"size:64;index" json:"namespace"`
	Type       SubjectType `gorm:"size:32;not null" json:"type"`
	ParentID   *string     `gorm:"size:64" json:"parent_id"`     // 仅用于org_node
	Members    JSONB       `gorm:"type:jsonb" json:"members"`    // 用户组：成员列表；组织节点：动态规则
	Attributes JSONB       `gorm:"type:jsonb" json:"attributes"` // 扩展属性
}

func (Subject) TableName() string {
	return "subjects"
}
