package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *AuthLoginHandler) handle(ctx context.Context, req AuthLoginRequest) (*AuthLoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "username or password is empty")
	}

	return &AuthLoginResponse{}, nil
}

func (h *AuthLoginHandler) GinHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req AuthLoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, GenBaseRespWithErr(err))
			return
		}

		resp, err := h.handle(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, GenBaseRespWithErr(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
