package api

import (
	"context"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/gin-gonic/gin"
)

type ListDocumentsResponse struct {
	BaseResponse
	Documents []dto.DocumentResponse `json:"documents"`
}

type ListDocumentsRequest struct {
	UserExternalID         string `json:"user_external_id"`
	RootDocumentExternalID string `json:"root_document_external_id"`
}

func (h *ListDocumentsHandler) handle(ctx context.Context, req Request) (Response, error) {

	return &ListDocumentsResponse{}, nil
}

func (h *ListDocumentsHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(ListDocumentsHandler{}), h.handle)
}
