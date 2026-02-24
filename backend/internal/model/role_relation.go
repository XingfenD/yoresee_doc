package model

type SubjectType int64

const (
	SubjectType_User      SubjectType = 0
	SubjectType_UserGroup SubjectType = 1
	SubjectType_OrgNode   SubjectType = 2
	SubjectType_AllUser   SubjectType = 3
)

type RoleRelation struct {
	ID          int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID      int64       `gorm:"not null;index:idx_role_subject" json:"role_id"`
	SubjectType SubjectType `gorm:"size:32;not null;index:idx_role_subject" json:"subject_type"`
	SubjectID   int64       `gorm:"not null;index:idx_role_subject" json:"subject_id"`
}

func (RoleRelation) TableName() string {
	return "role_relations"
}
