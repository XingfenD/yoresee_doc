package model

import (
	"time"
)

type Invitation struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code      string     `gorm:"size:32;unique;not null" json:"code"`
	CreatedBy int64      `gorm:"not null" json:"created_by"`
	IsUsed    bool       `gorm:"default:false" json:"is_used"`
	UsedAt    *time.Time `json:"used_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (Invitation) TableName() string {
	return "invitations"
}
