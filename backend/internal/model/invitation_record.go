package model

import "time"

type InvitationRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code         string    `gorm:"size:32;index;not null" json:"code"`
	UsedByUserID *int64    `gorm:"index" json:"used_by_user_id"`
	UsedBy       string    `gorm:"size:100" json:"used_by"`
	Status       string    `gorm:"size:20;index" json:"status"`
	UsedAt       time.Time `json:"used_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (InvitationRecord) TableName() string {
	return "invitation_records"
}
