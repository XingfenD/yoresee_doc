package model

type MembershipType int64

const (
	MembershipTypeUserGroup MembershipType = 1 // equals to SubjectTypeUserGroup
	MembershipTypeOrgNode   MembershipType = 2 // equals to SubjectTypeOrgNode
)

type MembershipRelation struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Type         MembershipType `gorm:"not null;index" json:"type"` // 1: UserGroup, 2: OrgNode
	UserID       int64          `gorm:"not null;index" json:"user_id"`
	MembershipID int64          `gorm:"not null;index" json:"membership_id"` // user_group_id or org_node_id
	User         User
}

// query all user belong to membership with type and member_id

func (MembershipRelation) TableName() string {
	return "membership"
}
