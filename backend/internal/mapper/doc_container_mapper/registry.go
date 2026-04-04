package doc_container_mapper

import (
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

const defaultContainerName = "own"

// Mapper defines a document container type mapping across transport(dto/proto) and model layers.
type Mapper interface {
	Name() string
	ProtoType() pb.CreateDocumentContainerType
	DTOType() dto.ContainerType
	ModelType() model.ContainerType
	RequiresKnowledgeBaseID() bool
}

type registry struct {
	byName  map[string]Mapper
	byProto map[pb.CreateDocumentContainerType]Mapper
	byDTO   map[dto.ContainerType]Mapper
	byModel map[model.ContainerType]Mapper
}

var globalRegistry = &registry{
	byName:  make(map[string]Mapper),
	byProto: make(map[pb.CreateDocumentContainerType]Mapper),
	byDTO:   make(map[dto.ContainerType]Mapper),
	byModel: make(map[model.ContainerType]Mapper),
}

func RegisterMapper(mapper Mapper) {
	if err := globalRegistry.register(mapper); err != nil {
		panic(err)
	}
}

func (r *registry) register(mapper Mapper) error {
	if mapper == nil {
		return fmt.Errorf("document container mapper is nil")
	}
	name := normalizeName(mapper.Name())
	if name == "" {
		return fmt.Errorf("document container mapper name is empty")
	}

	protoType := mapper.ProtoType()
	dtoType := mapper.DTOType()
	modelType := mapper.ModelType()
	if protoType == pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_UNSPECIFIED || dtoType == "" || modelType == "" {
		return fmt.Errorf("document container mapper %q has invalid mapping", name)
	}

	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("document container mapper %q already registered", name)
	}
	if _, exists := r.byProto[protoType]; exists {
		return fmt.Errorf("document container proto mapping %v already registered", protoType)
	}
	if _, exists := r.byDTO[dtoType]; exists {
		return fmt.Errorf("document container dto mapping %q already registered", dtoType)
	}
	if _, exists := r.byModel[modelType]; exists {
		return fmt.Errorf("document container model mapping %q already registered", modelType)
	}

	r.byName[name] = mapper
	r.byProto[protoType] = mapper
	r.byDTO[dtoType] = mapper
	r.byModel[modelType] = mapper
	return nil
}

func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (r *registry) findByName(name string) Mapper {
	return r.byName[normalizeName(name)]
}

func (r *registry) findByProto(t pb.CreateDocumentContainerType) Mapper {
	return r.byProto[t]
}

func (r *registry) findByDTO(t dto.ContainerType) Mapper {
	if mapper, ok := r.byDTO[t]; ok {
		return mapper
	}
	return r.byName[normalizeName(string(t))]
}

func (r *registry) findByModel(t model.ContainerType) Mapper {
	if mapper, ok := r.byModel[t]; ok {
		return mapper
	}
	return r.byName[normalizeName(string(t))]
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

func DefaultDTOType() dto.ContainerType {
	mapper := globalRegistry.defaultMapper()
	if mapper == nil {
		return dto.ContainerType(defaultContainerName)
	}
	return mapper.DTOType()
}

func IsSupportedDTOType(t dto.ContainerType) bool {
	return globalRegistry.findByDTO(t) != nil
}

func IsSupportedProtoType(t pb.CreateDocumentContainerType) bool {
	return globalRegistry.findByProto(t) != nil
}

func FromProtoType(t pb.CreateDocumentContainerType) (dto.ContainerType, bool) {
	if mapper := globalRegistry.findByProto(t); mapper != nil {
		return mapper.DTOType(), true
	}
	return "", false
}

func ToProtoType(t dto.ContainerType) pb.CreateDocumentContainerType {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.ProtoType()
	}
	return pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_UNSPECIFIED
}

func ToModelType(t dto.ContainerType) model.ContainerType {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.ModelType()
	}
	return model.ContainerType(DefaultDTOType())
}

func FromModelType(t model.ContainerType) dto.ContainerType {
	if mapper := globalRegistry.findByModel(t); mapper != nil {
		return mapper.DTOType()
	}
	normalized := normalizeName(string(t))
	if normalized != "" {
		return dto.ContainerType(normalized)
	}
	return DefaultDTOType()
}

func RequiresKnowledgeBaseID(t dto.ContainerType) bool {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.RequiresKnowledgeBaseID()
	}
	return false
}
