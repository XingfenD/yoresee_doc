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

	// userExternalID, _ := ctx.Value("user_external_id").(string)

	// allowed, err := service.DocumentSvc.CheckDocumentPermission(0, 0, "default", "read")
	// if err != nil {
	// 	return nil, err
	// }
	// if !allowed {
	// 	return nil, status.StatusPermissionDenied_DocumentRead
	// }

	document, content, err := service.DocumentSvc.GetDocumentWithContent(getDocReq.DocumentExternalID)
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
