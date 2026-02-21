package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *TestPostHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	testPostReq, ok := req.(TestPostRequest)
	if !ok {
		return nil, status.StatusParamError
	}
	if testPostReq.Message == "error" {
		return nil, status.StatusParamError
	}
	return &TestPostResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Message:      testPostReq.Message,
	}, nil
}

func (h *TestPostHandler) GinHandle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TestPostRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, GenBaseRespWithErrAndCtx(ctx, err))
			return
		}

		resp, err := h.handle(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, GenBaseRespWithErrAndCtx(ctx, err))
			return
		}

		if testResp, ok := resp.(*TestPostResponse); ok {
			testResp.Message = i18n.Translate(ctx, testResp.Message)
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
