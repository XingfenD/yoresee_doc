package model

import (
	"time"
)

type DocumentCollaborator struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentID int64     `gorm:"not null;index" json:"document_id"`
	UserID     int64     `gorm:"not null;index" json:"user_id"`
	Permission string    `gorm:"size:20;not null" json:"permission"` // read, write, admin
	CreatedAt  time.Time `json:"created_at"`
}

func (DocumentCollaborator) TableName() string {
	return "document_collaborators"
}
