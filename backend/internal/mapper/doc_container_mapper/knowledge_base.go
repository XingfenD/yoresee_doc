package doc_container_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type knowledgeBaseAdapter struct{}

func (knowledgeBaseAdapter) Name() string {
	return string(dto.ContainerType_KnowledgeBase)
}

func (knowledgeBaseAdapter) ProtoType() pb.CreateDocumentContainerType {
	return pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_KNOWLEDGE_BASE
}

func (knowledgeBaseAdapter) DTOType() dto.ContainerType {
	return dto.ContainerType_KnowledgeBase
}

func (knowledgeBaseAdapter) ModelType() model.ContainerType {
	return model.ContainerType_KnowledgeBase
}

func (knowledgeBaseAdapter) RequiresKnowledgeBaseID() bool {
	return true
}

func init() { RegisterMapper(knowledgeBaseAdapter{}) }
