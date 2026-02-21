package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ExternalID     string         `gorm:"size:100;unique" json:"external_id"`
	Username       string         `gorm:"size:50;unique;not null" json:"username"`
	Email          string         `gorm:"size:100;unique;not null" json:"email"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	Nickname       string         `gorm:"size:50" json:"nickname"`
	Avatar         string         `gorm:"size:255" json:"avatar"`
	RoleID         int64          `gorm:"not null" json:"role_id"`
	Role           Role           `gorm:"foreignKey:RoleID" json:"role"`
	Status         int            `gorm:"default:1" json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	InvitationCode *string        `gorm:"size:32" json:"invitation_code"`
}

func (User) TableName() string {
	return "users"
}
