package model

import (
	"time"
)

type Invitation struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Code       string     `gorm:"size:32;unique;not null" json:"code"`
	CreatedBy  int64      `gorm:"not null" json:"created_by"`
	UsedCnt    int64      `gorm:"default:0" json:"used_cnt"`
	MaxUsedCnt *int64     `gorm:"default:1" json:"max_used_cnt"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
	Disabled   bool       `gorm:"index,default:false" json:"disabled_at"`
}

func (Invitation) TableName() string {
	return "invitations"
}
