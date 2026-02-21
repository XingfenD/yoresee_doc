package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *AuthRegisterHandler) handle(ctx context.Context, req Request) (Response, error) {
	authRegisterReq := req.(AuthRegisterRequest)
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
	return func(ctx *gin.Context) {
		var req AuthRegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, GenBaseRespWithErrAndCtx(ctx, err))
			return
		}

		resp, err := h.handle(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, GenBaseRespWithErrAndCtx(ctx, err))
			return
		}

		if registerResp, ok := resp.(*AuthRegisterResponse); ok {
			registerResp.Message = i18n.Translate(ctx, registerResp.Message)
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
