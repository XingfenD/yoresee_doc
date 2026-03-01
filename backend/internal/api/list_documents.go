package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/types"
	"github.com/gin-gonic/gin"
)

type ListDocumentsResponse struct {
	api_base.BaseResponse
	Documents []dto.DocumentResponse `json:"documents"`
}

type ListDocumentsOptions struct {
	IncludeChildren bool `json:"include_children"`
	Recursive       bool `json:"recursive"`
	Depth           int  `json:"depth"`
}

type ListDocumentsArgs_OrderBy string

const (
	ListDocumentsArgs_OrderBy_CreatedAt ListDocumentsArgs_OrderBy = "created_at"
	ListDocumentsArgs_OrderBy_UpdatedAt ListDocumentsArgs_OrderBy = "updated_at"
)

type ListDocumentsRequest struct {
	UserExternalID         *string `json:"user_external_id" form:"user_external_id"`
	RootDocumentExternalID *string `json:"root_document_external_id" form:"root_document_external_id"`

	TitleKeyword *string  `json:"title_keyword,omitempty" form:"title_keyword"`
	Type         *string  `json:"type,omitempty" form:"type"`
	Status       *int     `json:"status,omitempty" form:"status"`
	Tags         []string `json:"tags,omitempty" form:"tags"`

	CreateTimeRange *types.TimeRange `json:"create_time_range,omitempty" form:"create_time_range"`
	UpdateTimeRange *types.TimeRange `json:"update_time_range,omitempty" form:"update_time_range"`

	OrderBy   *ListDocumentsArgs_OrderBy `json:"order_by" form:"order_by"`
	OrderDesc *bool                      `json:"order_desc" form:"order_desc"`

	Page     int `json:"page,omitempty" form:"page"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`

	Options *ListDocumentsOptions `json:"options,omitempty" form:"options"`
}

func (r *ListDocumentsRequest) BuildServiceReq() *service.ListDocumentsByExternalReq {
	if r == nil {
		return nil
	}

	sortArgs := service.SortArgs{
		Field: "created_at",
		Desc:  true,
	}
	if r.OrderBy != nil {
		switch *r.OrderBy {
		case ListDocumentsArgs_OrderBy_CreatedAt:
			sortArgs.Field = "created_at"
		case ListDocumentsArgs_OrderBy_UpdatedAt:
			sortArgs.Field = "updated_at"
		}
	}
	if r.OrderDesc != nil {
		sortArgs.Desc = *r.OrderDesc
	}

	filterArgs := &service.DocumentsListFilterArgs{
		TitleKeyword: r.TitleKeyword,
		DocType:      r.Type,
		Status:       r.Status,
		Tags:         r.Tags,
	}
	if r.CreateTimeRange != nil {
		filterArgs.CreateTimeRangeStart = r.CreateTimeRange.Start
		filterArgs.CreateTimeRangeEnd = r.CreateTimeRange.End
	}
	if r.UpdateTimeRange != nil {
		filterArgs.UpdateTimeRangeStart = r.UpdateTimeRange.Start
		filterArgs.UpdateTimeRangeEnd = r.UpdateTimeRange.End
	}
	req := &service.ListDocumentsByExternalReq{
		ExternalArgs: &service.DocumentsListExternalArgs{
			UserExternalID:         r.UserExternalID,
			RootDocumentExternalID: r.RootDocumentExternalID,
		},
		FilterArgs: filterArgs,
		SortArgs:   sortArgs,
		Pagination: service.Pagination{
			Page:     r.Page,
			PageSize: r.PageSize,
		},
		Options: r.Options.ToServiceOptions(),
	}
	return req
}

func (h *ListDocumentsHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	listDocsReq, ok := req.(*ListDocumentsRequest)
	if !ok {
		return nil, status.StatusParamError
	}

	responseDocsPtr, _, err := service.DocumentSvc.ListDocumentsWithChildrenByExternal(
		listDocsReq.BuildServiceReq(),
	)
	if err != nil {
		return nil, err
	}

	responseDocs := make([]dto.DocumentResponse, len(responseDocsPtr))
	for i, doc := range responseDocsPtr {
		responseDocs[i] = *doc
	}

	return &ListDocumentsResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Documents:    responseDocs,
	}, nil
}

func (apiOptions *ListDocumentsOptions) ToServiceOptions() *service.ListDocumentsOptions {
	if apiOptions == nil {
		return nil
	}
	return &service.ListDocumentsOptions{
		IncludeChildren: apiOptions.IncludeChildren,
		Recursive:       apiOptions.Recursive,
		Depth:           apiOptions.Depth,
	}
}

func (h *ListDocumentsHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(ListDocumentsRequest{}), h.handle)
}
