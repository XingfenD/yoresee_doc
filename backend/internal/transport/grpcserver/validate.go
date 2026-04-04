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
		if !doc_container_mapper.IsSupportedProtoType(*req.MoveToContainer) {
			return status.GenErrWithCustomMsg(status.StatusParamError, "invalid move_to_container")
		}
		if *req.MoveToContainer == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_OWN && req.KnowledgeBaseExternalId != nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId not nil")
		}
		if *req.MoveToContainer == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_KNOWLEDGE_BASE &&
			req.KnowledgeBaseExternalId == nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalId is nil")
		}
	}
	if req.MoveToContainer == nil && req.KnowledgeBaseExternalId != nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "MoveToContainer is nil")
	}
	return nil
}
