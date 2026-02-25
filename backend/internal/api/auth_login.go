package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	api_base.BaseResponse
	Token string           `json:"token"`
	User  dto.UserResponse `json:"user"`
}

func (h *AuthLoginHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	authLoginReq, ok := req.(*AuthLoginRequest)
	if !ok {
		return nil, status.StatusParamError
	}
	if authLoginReq.Email == "" || authLoginReq.Password == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "email or password is empty")
	}

	token, user, err := service.AuthSvc.Login(authLoginReq.Email, authLoginReq.Password)
	if err != nil {
		return nil, err
	}

	return &AuthLoginResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Token:        token,
		User:         *user,
	}, nil
}

func (h *AuthLoginHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(AuthLoginRequest{}), h.handle)
}
