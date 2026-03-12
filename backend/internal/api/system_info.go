package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type SystemInfoRequest struct{}

type SystemInfoResponse struct {
	api_base.BaseResponse
	SystemName         string `json:"system_name"`
	SystemRegisterMode string `json:"system_register_mode"`
}

func (h *SystemInfoHandler) handle(ctx context.Context, req api_base.Request) (resp api_base.Response, err error) {
	return &SystemInfoResponse{
		BaseResponse:       api_base.GenBaseRespWithErr(status.StatusSuccess),
		SystemName:         config.GlobalConfig.Backend.SystemName,
		SystemRegisterMode: service.ConfigSvc.GetSystemRegisterMode(ctx),
	}, nil
}

func (h *SystemInfoHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(SystemInfoRequest{}), h.handle)
}
