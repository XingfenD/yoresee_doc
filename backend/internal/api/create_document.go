package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/gin-gonic/gin"
)

type CreateDocumentContainerType string

const (
	CreateDocumentContainerType_KnowledgeBase CreateDocumentContainerType = "knowledge_base"
	CreateDocumentContainerType_Own           CreateDocumentContainerType = "own"
)

type CreateDocumentRequest struct {
	Content                 string                      `json:"content" form:"content"`
	Type                    string                      `json:"type" form:"type"`
	ContainerType           CreateDocumentContainerType `json:"container_type" form:"container_type"`
	KnowledgeBaseExternalID *string                     `json:"knowledge_base_external_id" form:"knowledge_base_external_id"`
	OwnerExternalID         *string                     `json:"owner_external_id" form:"owner_external_id"`
}

type CreateDocumentResponse struct {
	ExternalID string `json:"external_id"`
}

func (h *CreateDocumentHandler) validate(req *CreateDocumentRequest) error {
	if req == nil {
		return status.StatusParamError
	}

	if req.ContainerType == CreateDocumentContainerType_KnowledgeBase && req.KnowledgeBaseExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "nil KnowledgeBaseExternalID when ContainerType is knowledge_base")
	}

	if req.ContainerType == CreateDocumentContainerType_Own && req.OwnerExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "nil OwnerExternalID when ContainerType is own")
	}

	return nil
}

func (h *CreateDocumentHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	createDocumentReq := req.(*CreateDocumentRequest)
	if err := h.validate(createDocumentReq); err != nil {
		return nil, err
	}

	resp := &CreateDocumentResponse{}
	return resp, nil
}

func (h *CreateDocumentHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(CreateDocumentRequest{}), h.handle)
}
