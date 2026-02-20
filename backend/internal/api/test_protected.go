package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *TestProtectedHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &TestProtectedResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Message:      "protected",
	}, nil
}

func (h *TestProtectedHandler) GinHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := h.handle(ctx, TestProtectedRequest{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, GenBaseRespWithErr(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}
