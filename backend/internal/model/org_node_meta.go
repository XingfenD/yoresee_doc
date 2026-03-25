package model

type OrgNodeMeta struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`                 // org_node_id
	ExternalID  string `gorm:"not null;index;unique" json:"external_id"` // external org_node_id
	ParentID    int64  `gorm:"not null;index" json:"parent_id"`          // org_node_id, root node is 0
	Name        string `gorm:"not null;index" json:"name"`
	Path        string `gorm:"type:ltree;not null;index" json:"path"` // format: n<ID>.n<ID>...
	Description string `gorm:"type:text;index" json:"description"`
	CreatorID   int64  `gorm:"not null;index" json:"creator_id"` // user_id
}

func (OrgNodeMeta) TableName() string {
	return "org_node_meta"
}
