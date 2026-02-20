package model

import (
	"time"
)

type Role struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:50;unique;not null" json:"name"`
	Code        string    `gorm:"size:50;unique;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	IsSystem    bool      `gorm:"default:false" json:"is_system"`
	Permissions []string  `gorm:"serializer:json" json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}
