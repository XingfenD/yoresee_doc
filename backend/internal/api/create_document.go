package api

import (
	"context"
	"reflect"

	api_base "github.com/XingfenD/yoresee_doc/internal/api/base"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/gin-gonic/gin"
)

type CreateDocumentContainerType string

const (
	CreateDocumentContainerType_KnowledgeBase CreateDocumentContainerType = "knowledge_base"
	CreateDocumentContainerType_Own           CreateDocumentContainerType = "own"
)

type CreateDocumentRequest struct {
	Title                   string                      `json:"title" form:"title"`
	Type                    string                      `json:"type" form:"type"`
	ContainerType           CreateDocumentContainerType `json:"container_type" form:"container_type"`
	KnowledgeBaseExternalID *string                     `json:"knowledge_base_external_id" form:"knowledge_base_external_id"`
	ParentExternalID        *string                     `json:"parent_external_id" form:"parent_external_id"`
}

type CreateDocumentResponse struct {
	ExternalID string `json:"external_id"`
}

func (req *CreateDocumentRequest) BuildDTOReq(userExternalID string) *dto.CreateDocumentReq {
	if req == nil {
		return nil
	}

	dtoReq := &dto.CreateDocumentReq{
		Title:             req.Title,
		Type:              dto.DocumentType(req.Type),
		CreatorExternalID: utils.Of(userExternalID),
	}
	switch req.ContainerType {
	case CreateDocumentContainerType_KnowledgeBase:
		dtoReq.CreateAsOwnDoc = false
		dtoReq.KnowledgeExternalID = req.KnowledgeBaseExternalID
	case CreateDocumentContainerType_Own:
		dtoReq.CreateAsOwnDoc = true
	}
	return dtoReq
}
func (h *CreateDocumentHandler) handle(ctx context.Context, req api_base.Request) (api_base.Response, error) {
	createDocumentReq := req.(*CreateDocumentRequest)
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok {
		return nil, status.GenErrWithCustomMsg(status.StatusParamError, "user not logged in")
	}
	dtoReq := createDocumentReq.BuildDTOReq(userExternalID)
	resp, err := service.DocumentSvc.Create(dtoReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *CreateDocumentHandler) GinHandle() gin.HandlerFunc {
	baseHandler := &api_base.BaseHandler{}
	return baseHandler.GinHandle(reflect.TypeOf(CreateDocumentRequest{}), h.handle)
}
