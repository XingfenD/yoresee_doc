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

type AuthRegisterRequest struct {
	Username       string  `json:"username"`
	Password       string  `json:"password"`
	Email          string  `json:"email"`
	InvitationCode *string `json:"invitation_code"`
}

type AuthRegisterResponse struct {
	api_base.BaseResponse
}

func (h *AuthRegisterHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	authRegisterReq, ok := req.(*AuthRegisterRequest)
	if !ok {
		return nil, status.StatusParamError
	}
	if authRegisterReq.Username == "" || authRegisterReq.Password == "" || authRegisterReq.Email == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "username, password or email is empty")
	}

	userCreate := &dto.UserCreate{
		Username:       authRegisterReq.Username,
		Email:          authRegisterReq.Email,
		Password:       authRegisterReq.Password,
		InvitationCode: authRegisterReq.InvitationCode,
	}

	err := service.AuthSvc.Register(userCreate)
	if err != nil {
		return nil, err
	}

	return &AuthRegisterResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
	}, nil
}

func (h *AuthRegisterHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(AuthRegisterRequest{}), h.handle)
}
