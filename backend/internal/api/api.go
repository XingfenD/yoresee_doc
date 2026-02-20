package api

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type API interface {
	GinHander() gin.HandlerFunc
	handle(ctx context.Context, req Request) (resp Response, err error)
}

type Request interface {
}

type Response interface {
}

type BaseRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GenBaseRespWithErr(err error) BaseResponse {
	status, ok := err.(*status.Status)
	if !ok {
		return BaseResponse{
			Code:    500,
			Message: "internal server error",
		}
	}
	return BaseResponse{
		Code:    status.Code,
		Message: status.Message,
	}
}
