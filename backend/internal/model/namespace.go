package model

import (
	"time"
)

type Namespace struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	OwnerID   string    `gorm:"size:64;not null" json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Namespace) TableName() string {
	return "namespaces"
}
