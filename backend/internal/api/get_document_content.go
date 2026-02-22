package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

func (h *GetDocumentContentHandler) handle(ctx context.Context, req Request) (resp Response, err error) {
	getDocReq, ok := req.(*GetDocumentContentRequest)
	if !ok {
		return nil, status.StatusParamError
	}

	userExternalID, _ := ctx.Value("user_external_id").(string)
	userID, err := service.UserSvc.GetIDByExternalID(userExternalID)
	if err != nil {
		return nil, err
	}

	document, err := service.DocumentSvc.GetDocumentByExternalID(getDocReq.DocumentExternalID)
	if err != nil {
		return nil, err
	}

	allowed, err := service.DocumentSvc.CheckDocumentPermission(userID, document.ID, "read")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, status.StatusPermissionDenied_DocumentRead
	}

	content, err := service.DocumentSvc.GetDocumentContent(document.ID)
	if err != nil {
		return nil, err
	}

	return &GetDocumentContentResponse{
		BaseResponse: GenBaseRespWithErr(status.StatusSuccess),
		Document:     service.DocumentSvc.ConvertToDocumentResponse(document),
		Content:      content,
	}, nil
}

func (h *GetDocumentContentHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(GetDocumentContentRequest{}), h.handle)
}
