package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *AuthRegisterHandler) handle(ctx context.Context, req Request) (Response, error) {
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
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
	}, nil
}

func (h *AuthRegisterHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(AuthRegisterRequest{}), h.handle)
}
