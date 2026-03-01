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

type GetDocumentContentRequest struct {
	DocumentExternalID string `json:"document_external_id" form:"document_external_id" uri:"documentExternalID"`
}

type GetDocumentContentResponse struct {
	api_base.BaseResponse
	Document *dto.DocumentResponse `json:"document"`
	Content  string                `json:"content"`
}

func (h *GetDocumentContentHandler) handle(ctx context.Context, req api_base.Request) (resp api_base.Response, err error) {
	getDocReq, ok := req.(*GetDocumentContentRequest)
	if !ok {
		return nil, status.StatusParamError
	}

	// userExternalID, _ := ctx.Value("user_external_id").(string)
	// userID, err := service.UserSvc.GetIDByExternalID(userExternalID)
	// if err != nil {
	// 	return nil, status.StatusUserNotFound
	// }

	document, err := service.DocumentSvc.GetDocumentByExternalID(getDocReq.DocumentExternalID)
	if err != nil {
		return nil, status.StatusDocumentNotFound
	}

	// allowed, err := service.DocumentSvc.CheckDocumentPermission(userID, document.ID, string(model.PermissionRead))
	// if err != nil {
	// 	return nil, err
	// }
	// if !allowed {
	// 	return nil, status.StatusPermissionDenied_DocumentRead
	// }

	content, err := service.DocumentSvc.GetDocumentContent(document.ID)
	if err != nil {
		return nil, err
	}

	return &GetDocumentContentResponse{
		BaseResponse: api_base.GenBaseRespWithErr(status.StatusSuccess),
		Document:     service.DocumentSvc.ConvertToDocumentResponse(document),
		Content:      content,
	}, nil
}

func (h *GetDocumentContentHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(GetDocumentContentRequest{}), h.handle)
}
