package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type TestPostRequest struct {
	Message string `json:"message"`
}

type TestPostResponse struct {
	api_base.BaseResponse
	Message string `json:"message"`
}

func (h *TestPostHandler) handle(ctx context.Context, req api_base.Request) (resp api_base.Response, err error) {
	testPostReq, ok := req.(*TestPostRequest)
	if !ok {
		return nil, status.StatusParamError
	}
	if testPostReq.Message == "error" {
		return nil, status.StatusParamError
	}
	return &TestPostResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Message:      testPostReq.Message,
	}, nil
}

func (h *TestPostHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(TestPostRequest{}), h.handle)
}
