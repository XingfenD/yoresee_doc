package template_container_mapper

import (
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

const defaultContainerName = "own"

type Mapper interface {
	Name() string
	ProtoType() pb.CreateTemplateContainer
	DTOType() dto.TemplateContainer
	Scope() string
	IsPublic() bool
	RequiresKnowledgeBaseID() bool
}

type registry struct {
	byName  map[string]Mapper
	byProto map[pb.CreateTemplateContainer]Mapper
	byDTO   map[dto.TemplateContainer]Mapper
}

var globalRegistry = &registry{
	byName:  make(map[string]Mapper),
	byProto: make(map[pb.CreateTemplateContainer]Mapper),
	byDTO:   make(map[dto.TemplateContainer]Mapper),
}

func RegisterMapper(mapper Mapper) {
	if err := globalRegistry.register(mapper); err != nil {
		panic(err)
	}
}

func (r *registry) register(mapper Mapper) error {
	if mapper == nil {
		return fmt.Errorf("template container mapper is nil")
	}
	name := normalizeName(mapper.Name())
	if name == "" {
		return fmt.Errorf("template container mapper name is empty")
	}

	protoType := mapper.ProtoType()
	dtoType := mapper.DTOType()
	if dtoType < dto.TemplateContainerOwn {
		return fmt.Errorf("template container mapper %q has invalid dto mapping", name)
	}
	if strings.TrimSpace(mapper.Scope()) == "" {
		return fmt.Errorf("template container mapper %q has empty scope mapping", name)
	}

	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("template container mapper %q already registered", name)
	}
	if _, exists := r.byProto[protoType]; exists {
		return fmt.Errorf("template container proto mapping %v already registered", protoType)
	}
	if _, exists := r.byDTO[dtoType]; exists {
		return fmt.Errorf("template container dto mapping %d already registered", dtoType)
	}

	r.byName[name] = mapper
	r.byProto[protoType] = mapper
	r.byDTO[dtoType] = mapper
	return nil
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (r *registry) findByName(name string) Mapper {
	return r.byName[normalizeName(name)]
}

func (r *registry) findByProto(t pb.CreateTemplateContainer) Mapper {
	return r.byProto[t]
}

func (r *registry) findByDTO(t dto.TemplateContainer) Mapper {
	return r.byDTO[t]
}

func (r *registry) defaultMapper() Mapper {
	if mapper := r.findByName(defaultContainerName); mapper != nil {
		return mapper
	}
	for _, mapper := range r.byName {
		return mapper
	}
	return nil
}

func DefaultDTOType() dto.TemplateContainer {
	mapper := globalRegistry.defaultMapper()
	if mapper == nil {
		return dto.TemplateContainerOwn
	}
	return mapper.DTOType()
}

func IsSupportedDTOType(t dto.TemplateContainer) bool {
	return globalRegistry.findByDTO(t) != nil
}

func FromProtoType(t pb.CreateTemplateContainer) dto.TemplateContainer {
	if mapper := globalRegistry.findByProto(t); mapper != nil {
		return mapper.DTOType()
	}
	return DefaultDTOType()
}

func ToProtoType(t dto.TemplateContainer) pb.CreateTemplateContainer {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.ProtoType()
	}
	return pb.CreateTemplateContainer_OWN_TEMPLATE
}

func ToScope(t dto.TemplateContainer) string {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.Scope()
	}
	if mapper := globalRegistry.defaultMapper(); mapper != nil {
		return mapper.Scope()
	}
	return "private"
}

func ToIsPublic(t dto.TemplateContainer) bool {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.IsPublic()
	}
	if mapper := globalRegistry.defaultMapper(); mapper != nil {
		return mapper.IsPublic()
	}
	return false
}

func RequiresKnowledgeBaseID(t dto.TemplateContainer) bool {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.RequiresKnowledgeBaseID()
	}
	return false
}
