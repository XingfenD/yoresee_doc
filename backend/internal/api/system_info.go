package api

import (
	"context"
	"net/http"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *SystemInfoHandler) handle(ctx context.Context, req Request) (resp Response, err error) {

	return &SystemInfoResponse{
		BaseResponse:       GenBaseRespWithErr(status.StatusSuccess),
		SystemName:         config.GlobalConfig.Backend.SystemName,
		SystemRegisterMode: service.ConfigSvc.GetSystemRegisterMode(),
	}, nil
}

func (h *SystemInfoHandler) GinHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := &SystemInfoRequest{}
		resp, err := h.handle(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenBaseRespWithErrAndCtx(c, err))
			return
		}

		if infoResp, ok := resp.(*SystemInfoResponse); ok {
			infoResp.Message = i18n.Translate(c, infoResp.Message)
		}
		c.JSON(http.StatusOK, resp)
	}
}
