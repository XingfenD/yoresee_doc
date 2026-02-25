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

type ListDocumentsArgs struct {
	OrderBy   ListDocumentsArgs_OrderBy `json:"order_by"`
	OrderDesc bool                      `json:"order_desc"`
}

type ListDocumentsRequest struct {
	UserExternalID         *string `json:"user_external_id"`
	RootDocumentExternalID *string `json:"root_document_external_id"`

	TitleKeyword *string  `json:"title_keyword,omitempty"`
	Type         *string  `json:"type,omitempty"`
	Status       *int     `json:"status,omitempty"`
	Tags         []string `json:"tags,omitempty"`

	CreateTimeRange *TimeRange `json:"create_time_range,omitempty"`
	UpdateTimeRange *TimeRange `json:"update_time_range,omitempty"`

	OrderBy *ListDocumentsArgs `json:"order_by,omitempty"`

	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`

	Options *ListDocumentsOptions `json:"options,omitempty"`
}

type TimeRange struct {
	Start *string `json:"start"`
	End   *string `json:"end"`
}

func (h *ListDocumentsHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	listDocsReq, ok := req.(*ListDocumentsRequest)
	if !ok {
		return nil, status.StatusParamError
	}

	// 处理排序
	sortField := "created_at" // 默认排序
	sortDesc := true          // 默认降序
	if listDocsReq.OrderBy != nil {
		switch listDocsReq.OrderBy.OrderBy {
		case ListDocumentsArgs_OrderBy_CreatedAt:
			sortField = "created_at"
		case ListDocumentsArgs_OrderBy_UpdatedAt:
			sortField = "updated_at"
		}
		sortDesc = listDocsReq.OrderBy.OrderDesc
	}

	// 转换为服务层选项
	serviceOptions := convertListDocumentsOptions(listDocsReq.Options)

	// 时间范围已为字符串格式，直接使用
	var createTimeStart, createTimeEnd, updateTimeRangeStart, updateTimeRangeEnd *string
	if listDocsReq.CreateTimeRange != nil {
		createTimeStart = listDocsReq.CreateTimeRange.Start
		createTimeEnd = listDocsReq.CreateTimeRange.End
	}
	if listDocsReq.UpdateTimeRange != nil {
		updateTimeRangeStart = listDocsReq.UpdateTimeRange.Start
		updateTimeRangeEnd = listDocsReq.UpdateTimeRange.End
	}

	// 使用服务层查询文档（包含子文档）
	responseDocsPtr, _, err := service.DocumentSvc.ListDocumentsWithChildrenByExternalID(
		listDocsReq.UserExternalID,
		listDocsReq.RootDocumentExternalID,
		listDocsReq.TitleKeyword,
		listDocsReq.Type,
		listDocsReq.Status,
		listDocsReq.Tags,
		createTimeStart,
		createTimeEnd,
		updateTimeRangeStart,
		updateTimeRangeEnd,
		sortField,
		sortDesc,
		listDocsReq.Page,
		listDocsReq.PageSize,
		serviceOptions,
	)
	if err != nil {
		return nil, err
	}

	// 将 []*dto.DocumentResponse 转换为 []dto.DocumentResponse
	responseDocs := make([]dto.DocumentResponse, len(responseDocsPtr))
	for i, doc := range responseDocsPtr {
		responseDocs[i] = *doc
	}

	return &ListDocumentsResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Documents:    responseDocs,
	}, nil
}

// 将子文档查询委托给服务层
// API 层仅负责将 API 结构转换为服务层结构
func convertListDocumentsOptions(apiOptions *ListDocumentsOptions) *service.ListDocumentsOptions {
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
