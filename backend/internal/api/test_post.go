package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *TestPostHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	testPostReq, ok := req.(*TestPostRequest)
	if !ok {
		return nil, status.StatusParamError
	}
	if testPostReq.Message == "error" {
		return nil, status.StatusParamError
	}
	return &TestPostResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Message:      testPostReq.Message,
	}, nil
}

func (h *TestPostHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(TestPostRequest{}), h.handle)
}
