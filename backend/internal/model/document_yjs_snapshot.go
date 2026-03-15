package model

import "time"

type DocumentYjsSnapshot struct {
	DocID     int64     `gorm:"primaryKey;column:doc_id" json:"doc_id"`
	YjsState  []byte    `gorm:"type:bytea;not null" json:"yjs_state"`
	Version   int64     `gorm:"not null" json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DocumentYjsSnapshot) TableName() string {
	return "documents_yjs_snapshot"
}
