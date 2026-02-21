package api

import "github.com/XingfenD/yoresee_doc/internal/dto"

var HealthHandlerImpl = &HealthHandler{}
var TestProtectedHandlerImpl = &TestProtectedHandler{}
var TestPostHandlerImpl = &TestPostHandler{}
var AuthRegisterHandlerImpl = &AuthRegisterHandler{}
var AuthLoginHandlerImpl = &AuthLoginHandler{}
var SystemInfoHandlerImpl = &SystemInfoHandler{}

type HealthHandler struct{}
type TestProtectedHandler struct{}
type TestPostHandler struct{}
type AuthRegisterHandler struct{}
type AuthLoginHandler struct{}
type SystemInfoHandler struct{}

type HealthResponse struct {
	BaseResponse
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Version   string `json:"version"`
}

type HealthRequest struct {
	BaseRequest
}

type TestProtectedResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type TestProtectedRequest struct {
	BaseRequest
}

type TestPostResponse struct {
	BaseResponse
	Message string `json:"message"`
}

type TestPostRequest struct {
	BaseRequest
	Message string `json:"message"`
}

type AuthRegisterResponse struct {
	BaseResponse
}

type AuthRegisterRequest struct {
	BaseRequest
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	Email          string  `json:"email"`
	InvitationCode *string `json:"invitation_code"`
}

type AuthLoginResponse struct {
	BaseResponse
	Token string           `json:"token"`
	User  dto.UserResponse `json:"user"`
}

type AuthLoginRequest struct {
	BaseRequest
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SystemInfoResponse struct {
	BaseResponse
	SystemName         string `json:"system_name"`
	SystemRegisterMode string `json:"system_register_mode"`
}

type SystemInfoRequest struct {
	BaseRequest
}
