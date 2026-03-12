package api

import (
	"context"
	"reflect"
	"time"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GetKnowledgeBaseRequset struct {
	KnowledgeBaseExternalID string `json:"knowledge_base_external_id" form:"knowledge_base_external_id" uri:"knowledgeBaseExternalID"`
	RecordRecentLog         bool   `json:"record_recent_log" form:"record_recent_log"`
	Page                    int    `json:"page,omitempty" form:"page"`
	PageSize                int    `json:"page_size,omitempty" form:"page_size"`
}

type GetKnowledgeBaseResponse struct {
	api_base.BaseResponse
	KnowledgeBase *dto.KnowledgeBaseResponse `json:"knowledge_base"`
	Documents     []*dto.DocumentResponse    `json:"documents"`
	TotalCount    int64                      `json:"total_count"`
}

func (h *GetKnowledgeBaseHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	getKnowledgeBaseReq := req.(*GetKnowledgeBaseRequset)

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user not logged in")
	}

	knowledgeBaseDTO, err := service.KnowledgeBaseSvc.GetByExternalID(&dto.KnowledgeBaseGetByExternalIDReq{
		KnowledgeBaseExternalID: getKnowledgeBaseReq.KnowledgeBaseExternalID,
	}).WithExtend().Exec()
	if err != nil {
		return nil, status.StatusKnowledgeBaseNotFound
	}

	// if !knowledgeBaseDTO.IsPublic {
	// 	if knowledgeBaseDTO.CreatorUserExternalID != userExternalID {
	// 		return nil, status.StatusPermissionDenied
	// 	}
	// 	// TODO: permission check
	// }

	if getKnowledgeBaseReq.RecordRecentLog {
		service.KnowledgeBaseSvc.CreateRecentKnowledgeBase(&dto.CreateRecentKnowledgeBaseRequest{
			UserExternalID:          userExternalID,
			KnowledgeBaseExternalID: getKnowledgeBaseReq.KnowledgeBaseExternalID,
			AssessTime:              time.Now(),
		})
	}

	documents, totalCount, err := service.DocumentSvc.ListDocumentsWithChildrenByExternal(
		&dto.ListDocumentsByExternalReq{
			ExternalArgs: &dto.DocumentsListExternalArgs{
				KnowledgeExternalID: &knowledgeBaseDTO.ExternalID,
			},
			Pagination: dto.Pagination{
				Page:     getKnowledgeBaseReq.Page,
				PageSize: getKnowledgeBaseReq.PageSize,
			},
			Options: &dto.RecursiveOptions{
				IncludeChildren: true,
				Recursive:       true,
			},
		},
	)
	if err != nil {
		logrus.Errorf("[GetKnowledgeBaseHandler] ListDocumentsWithChildrenByExternal Failed: error: %+v", err)
		return nil, status.StatusDocumentNotFound
	}

	response := &GetKnowledgeBaseResponse{
		BaseResponse:  api_base.GenBaseRespWithErr(status.StatusSuccess),
		KnowledgeBase: knowledgeBaseDTO,
		Documents:     documents,
		TotalCount:    totalCount,
	}

	return response, nil
}

func (h *GetKnowledgeBaseHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(GetKnowledgeBaseRequset{}), h.handle)
}
