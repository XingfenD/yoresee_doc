package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *AuthLoginHandler) handle(ctx context.Context, req Request) (Response, error) {
	authLoginReq := req.(AuthLoginRequest)
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
	return func(ctx *gin.Context) {
		var req AuthLoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, GenBaseRespWithErr(err))
			return
		}

		logrus.Infof("AuthLoginRequest: %+v", req)
		resp, err := h.handle(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, GenBaseRespWithErr(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
