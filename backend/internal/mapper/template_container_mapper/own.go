package template_container_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type ownAdapter struct{}

func (ownAdapter) Name() string {
	return "own"
}

func (ownAdapter) ProtoType() pb.CreateTemplateContainer {
	return pb.CreateTemplateContainer_OWN_TEMPLATE
}

func (ownAdapter) DTOType() dto.TemplateContainer {
	return dto.TemplateContainerOwn
}

func (ownAdapter) Scope() string {
	return "private"
}

func (ownAdapter) IsPublic() bool {
	return false
}

func (ownAdapter) RequiresKnowledgeBaseID() bool {
	return false
}

func init() { RegisterMapper(ownAdapter{}) }
