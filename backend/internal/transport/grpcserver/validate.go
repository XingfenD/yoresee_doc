package grpcserver

import (
	doc_container_mapper "github.com/XingfenD/yoresee_doc/internal/mapper/doc_container_mapper"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

func validateUpdateDocumentRequest(req *pb.UpdateDocumentRequest) error {
	if req == nil {
		return status.StatusParamError
	}
	if req.Content == nil && req.KnowledgeBaseExternalId == nil && req.ParentExternalId == nil && req.Title == nil && req.MoveToContainer == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update fields")
	}
	if req.MoveToContainer != nil {
		containerType, ok := doc_container_mapper.FromProtoType(*req.MoveToContainer)
		if !ok {
			return status.GenErrWithCustomMsg(status.StatusParamError, "invalid move_to_container")
		}
		requiresKnowledgeBaseID := doc_container_mapper.RequiresKnowledgeBaseID(containerType)
		if !requiresKnowledgeBaseID && req.KnowledgeBaseExternalId != nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId not nil")
		}
		if requiresKnowledgeBaseID && req.KnowledgeBaseExternalId == nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId is nil")
		}
	}
	if req.MoveToContainer == nil && req.KnowledgeBaseExternalId != nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "MoveToContainer is nil")
	}
	return nil
}
