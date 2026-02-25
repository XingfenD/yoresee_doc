package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type SystemInfoRequest struct{}

type SystemInfoResponse struct {
	BaseResponse
	SystemName         string `json:"system_name"`
	SystemRegisterMode string `json:"system_register_mode"`
}

func (h *SystemInfoHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &SystemInfoResponse{
		BaseResponse:       GenBaseRespWithErr(status.StatusSuccess),
		SystemName:         config.GlobalConfig.Backend.SystemName,
		SystemRegisterMode: service.ConfigSvc.GetSystemRegisterMode(),
	}, nil
}

func (h *SystemInfoHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(SystemInfoRequest{}), h.handle)
}
