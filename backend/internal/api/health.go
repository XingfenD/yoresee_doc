package api

import (
	"context"
	"reflect"
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
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(HealthRequest{}), h.handle)
}
