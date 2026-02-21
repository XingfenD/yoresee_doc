package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *SystemInfoHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &SystemInfoResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		SystemName:   config.GlobalConfig.Backend.SystemName,
	}, nil
}

func (h *SystemInfoHandler) GinHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &SystemInfoRequest{}
		resp, err := h.handle(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenBaseRespWithErr(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
