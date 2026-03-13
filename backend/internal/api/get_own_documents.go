package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type GetOwnDocumentsRequest struct {
	Page     int `json:"page,omitempty" form:"page"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`
}

type GetOwnDocumentsResponse struct {
	api_base.BaseResponse
	Documents  []*dto.DocumentResponse `json:"documents"`
	TotalCount int64                   `json:"total_count"`
}

func (h *GetOwnDocumentsHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	getOwnDocsReq := req.(*GetOwnDocumentsRequest)

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user not logged in")
	}

	queryReq := &dto.ListDocumentsByExternalReq{
		ExternalArgs: &dto.DocumentsListExternalArgs{
			UserExternalID: &userExternalID,
			ListOwnDoc:     true,
		},
		Pagination: dto.Pagination{
			Page:     getOwnDocsReq.Page,
			PageSize: getOwnDocsReq.PageSize,
		},
	}

	docs, count, err := service.DocumentSvc.ListDocumentsWithChildrenByExternal(queryReq)
	if err != nil {
		return nil, err
	}

	return &GetOwnDocumentsResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Documents:    docs,
		TotalCount:   count,
	}, nil

}

func (h *GetOwnDocumentsHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(GetOwnDocumentsRequest{}), h.handle)
}
