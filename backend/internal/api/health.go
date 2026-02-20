package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *HealthHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &HealthResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Status:    "healthy",
		Version:   "1.0.0",
	}, nil
}

func (h *HealthHandler) GinHander() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &HealthRequest{}
		resp, err := h.handle(ctx, req)
		if err != nil {
			c.JSON(200, GenBaseRespWithErr(err))
			return
		}
		c.JSON(200, resp)
	}
}
