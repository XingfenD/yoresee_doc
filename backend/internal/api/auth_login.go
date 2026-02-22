package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *AuthLoginHandler) handle(ctx context.Context, req Request) (Response, error) {
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
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Token:        token,
		User:         *user,
	}, nil
}

func (h *AuthLoginHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(AuthLoginRequest{}), h.handle)
}
