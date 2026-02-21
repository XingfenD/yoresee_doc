package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/gin-gonic/gin"
)

func (h *SystemInfoHandler) handle(ctx context.Context, req Request) (resp Response, err error) {

	return &SystemInfoResponse{
		BaseResponse:       GenBaseRespWithErr(status.StatusSuccess),
		SystemName:         config.GlobalConfig.Backend.SystemName,
		SystemRegisterMode: getSystemRegisterMode(),
	}, nil
}

func getSystemRegisterMode() string {
	registerMode, err := service.ConfigSvc.Get(utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	))
	if err != nil {
		return constant.RegisterMode_Invite
	}
	return registerMode
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
