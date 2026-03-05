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
}

type GetKnowledgeBaseResponse struct {
	api_base.BaseResponse
	KnowledgeBase *dto.KnowledgeBaseResponse `json:"knowledge_base"`
	Documents     []*dto.DocumentResponse    `json:"documents"`
}

func (h *GetKnowledgeBaseHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	getKnowledgeBaseReq := req.(*GetKnowledgeBaseRequset)

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user not logged in")
	}

	knowledgeBaseDTO, err := service.KnowledgeBaseSvc.GetByExternalID(&service.KnowledgeBaseGetByExternalIDReq{
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

	documents, _, err := service.DocumentSvc.ListDocumentsWithChildrenByExternal(
		&service.ListDocumentsByExternalReq{
			ExternalArgs: &service.DocumentsListExternalArgs{
				KnowledgeExternalID: &knowledgeBaseDTO.ExternalID,
			},
			Options: &service.ListDocumentsOptions{
				IncludeChildren: true,
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
	}

	return response, nil
}

func (h *GetKnowledgeBaseHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(GetKnowledgeBaseRequset{}), h.handle)
}
