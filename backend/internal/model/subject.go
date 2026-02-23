package model

// SubjectType 主体类型枚举
type SubjectType int64

const (
	SubjectTypeUser      SubjectType = 0 // user
	SubjectTypeUserGroup SubjectType = 1 // user_group
	SubjectTypeOrgNode   SubjectType = 2 // org_node
)

type Subject struct {
	ID   int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Type SubjectType `gorm:"size:32;not null" json:"type"`
}

func (Subject) TableName() string {
	return "subjects"
}
