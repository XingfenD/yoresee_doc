package dto

import (
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

type UserBase struct {
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserCreate struct {
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	InvitationCode *string `json:"invitation_code"`
}

type UserUpdate struct {
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	Nickname       string  `json:"nickname"`
	Avatar         string  `json:"avatar"`
	Status         int     `json:"status"`
	InvitationCode *string `json:"invitation_code"`
}

type UserResponse struct {
	UserBase
	InvitationCode *string `json:"invitation_code"`
}

func NewUserResponseFromModel(user *model.User) *UserResponse {
	return &UserResponse{
		UserBase: UserBase{
			ExternalID: user.ExternalID,
			Username:   user.Username,
			Email:      user.Email,
			Nickname:   user.Nickname,
			Avatar:     user.Avatar,
			Status:     user.Status,
			CreatedAt:  user.CreatedAt,
			UpdatedAt:  user.UpdatedAt,
		},
		InvitationCode: user.InvitationCode,
	}
}
