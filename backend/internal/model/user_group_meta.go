package model

type UserGroupMeta struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`          // user_group_id
	ExternalID  string `gorm:"not null;index" json:"external_id"` // external user_group_id
	Name        string `gorm:"not null;index" json:"name"`
	CreatorID   int64  `gorm:"not null;index" json:"creator_id"` // user_id
	Description string `gorm:"type:text;index" json:"description"`
}

func (UserGroupMeta) TableName() string {
	return "user_group_meta"
}
