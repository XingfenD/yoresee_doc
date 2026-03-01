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

type ListKnowledgeBasesResponse struct {
	api_base.BaseResponse
	KnowledgeBases []*dto.KnowledgeBaseResponse `json:"knowledge_bases"`
}

type ListKnowledgeBasesOptions struct {
	IncludeChildren bool `json:"include_children"`
	Recursive       bool `json:"recursive"`
	Depth           int  `json:"depth"`
}

type ListKnowledgeBasesArgs_OrderBy string

const (
	ListKnowledgeBasesArgs_OrderBy_CreatedAt ListKnowledgeBasesArgs_OrderBy = "created_at"
	ListKnowledgeBasesArgs_OrderBy_UpdatedAt ListKnowledgeBasesArgs_OrderBy = "updated_at"
)

type ListKnowledgeBasesRequest struct {
	UserExternalID *string `json:"user_external_id" form:"user_external_id"`

	NameKeyword *string  `json:"name_keyword,omitempty" form:"name_keyword"`
	IsPublic    *bool    `json:"is_public,omitempty" form:"is_public"`
	Tags        []string `json:"tags,omitempty" form:"tags"`

	CreateTimeRange *types.TimeRange `json:"create_time_range,omitempty" form:"create_time_range"`
	UpdateTimeRange *types.TimeRange `json:"update_time_range,omitempty" form:"update_time_range"`

	OrderBy   *ListKnowledgeBasesArgs_OrderBy `json:"order_by" form:"order_by"`
	OrderDesc *bool                           `json:"order_desc" form:"order_desc"`

	Page     int `json:"page,omitempty" form:"page"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`

	// Options *ListKnowledgeBasesOptions `json:"options,omitempty"`
}

func (r *ListKnowledgeBasesRequest) BuildServiceReq() *service.KnowledgeBaseListByExternalReq {
	if r == nil {
		return nil
	}

	sortArgs := service.SortArgs{
		Field: "created_at",
		Desc:  true,
	}
	if r.OrderBy != nil {
		switch *r.OrderBy {
		case ListKnowledgeBasesArgs_OrderBy_CreatedAt:
			sortArgs.Field = "created_at"
		case ListKnowledgeBasesArgs_OrderBy_UpdatedAt:
			sortArgs.Field = "updated_at"
		}
	}
	if r.OrderDesc != nil {
		sortArgs.Desc = *r.OrderDesc
	}

	filterArgs := &service.KnowledgeBaseListFilterArgs{
		IsPublic: r.IsPublic,
	}

	req := &service.KnowledgeBaseListByExternalReq{
		CreatorExternalID: "",
		FilterArgs:        filterArgs,
		SortArgs:          sortArgs,
		Pagination: service.Pagination{
			Page:     r.Page,
			PageSize: r.PageSize,
		},
	}

	if r.UserExternalID != nil {
		req.CreatorExternalID = *r.UserExternalID
	}

	return req
}

func (h *ListKnowledgeBasesHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	listKBsReq, ok := req.(*ListKnowledgeBasesRequest)
	if !ok {
		return nil, status.StatusParamError
	}

	kbModels, err := service.KnowledgeBaseSvc.ListByExternal(
		listKBsReq.BuildServiceReq(),
	)
	if err != nil {
		return nil, err
	}

	return &ListKnowledgeBasesResponse{
		BaseResponse:   api_base.GenBaseRespWithErr(status.StatusSuccess),
		KnowledgeBases: kbModels,
	}, nil
}

func (h *ListKnowledgeBasesHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(ListKnowledgeBasesRequest{}), h.handle)
}
