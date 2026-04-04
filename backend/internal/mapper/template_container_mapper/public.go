package template_container_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type publicAdapter struct{}

func (publicAdapter) Name() string {
	return "public"
}

func (publicAdapter) ProtoType() pb.CreateTemplateContainer {
	return pb.CreateTemplateContainer_PUBLIC_TEMPLATE
}

func (publicAdapter) DTOType() dto.TemplateContainer {
	return dto.TemplateContainerPublic
}

func (publicAdapter) Scope() string {
	return "system"
}

func (publicAdapter) IsPublic() bool {
	return true
}

func (publicAdapter) RequiresKnowledgeBaseID() bool {
	return false
}

func init() { RegisterMapper(publicAdapter{}) }
