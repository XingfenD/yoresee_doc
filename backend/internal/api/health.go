package api

import (
	"context"
	"reflect"
	"time"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type HealthRequest struct{}

type HealthResponse struct {
	api_base.BaseResponse
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
	Version   string `json:"version"`
}

func (h *HealthHandler) handle(ctx context.Context, req api_base.Request) (resp api_base.Response, err error) {
	return &HealthResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Timestamp:    time.Now().Format(time.RFC3339),
		Status:       "healthy",
		Version:      "1.0.0",
	}, nil
}

func (h *HealthHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(HealthRequest{}), h.handle)
}
