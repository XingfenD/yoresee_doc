package template_container_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type knowledgeBaseAdapter struct{}

func (knowledgeBaseAdapter) Name() string {
	return "knowledge_base"
}

func (knowledgeBaseAdapter) ProtoType() pb.CreateTemplateContainer {
	return pb.CreateTemplateContainer_KNOWLEDGEBASE_TEMPLATE
}

func (knowledgeBaseAdapter) DTOType() dto.TemplateContainer {
	return dto.TemplateContainerKnowledgeBase
}

func (knowledgeBaseAdapter) Scope() string {
	return "knowledge_base"
}

func (knowledgeBaseAdapter) IsPublic() bool {
	return false
}

func (knowledgeBaseAdapter) RequiresKnowledgeBaseID() bool {
	return true
}

func init() { RegisterMapper(knowledgeBaseAdapter{}) }
