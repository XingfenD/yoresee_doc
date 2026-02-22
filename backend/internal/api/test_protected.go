package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *TestProtectedHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	return &TestProtectedResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Message:      "protected",
	}, nil
}

func (h *TestProtectedHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(TestProtectedRequest{}), h.handle)
}
