package api

import (
	"context"
	"net/http"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *HealthHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &HealthResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Timestamp:    time.Now().Format(time.RFC3339),
		Status:       "healthy",
		Version:      "1.0.0",
	}, nil
}

func (h *HealthHandler) GinHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &HealthRequest{}
		resp, err := h.handle(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenBaseRespWithErr(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
