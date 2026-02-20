package model

import "time"

type RegisterRecord struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         int64     `gorm:"not null;index" json:"user_id"`
	InvitationCode string    `gorm:"default:null" json:"invitation_code"`
	CreatedAt      time.Time `json:"created_at"`
}

func (RegisterRecord) TableName() string {
	return "register_records"
}
