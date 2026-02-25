package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type TestProtectedRequest struct{}

type TestProtectedResponse struct {
	api_base.BaseResponse
	Message string `json:"message"`
}

func (h *TestProtectedHandler) handle(ctx context.Context, req api_base.Request) (resp api_base.Response, err error) {
	return &TestProtectedResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Message:      "protected",
	}, nil
}

func (h *TestProtectedHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(TestProtectedRequest{}), h.handle)
}
