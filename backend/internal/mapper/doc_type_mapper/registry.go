package doc_type_mapper

import (
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

const defaultTypeName = "markdown"

// Mapper defines a document type mapping across transport(dto/proto) and model layers.
type Mapper interface {
	Name() string
	ProtoType() pb.DocumentType
	DTOType() dto.DocumentType
	ModelType() model.DocumentType
}

type registry struct {
	byName  map[string]Mapper
	byProto map[pb.DocumentType]Mapper
	byDTO   map[dto.DocumentType]Mapper
	byModel map[model.DocumentType]Mapper
}

var globalRegistry = &registry{
	byName:  make(map[string]Mapper),
	byProto: make(map[pb.DocumentType]Mapper),
	byDTO:   make(map[dto.DocumentType]Mapper),
	byModel: make(map[model.DocumentType]Mapper),
}

// RegisterMapper registers a document type mapper.
// Implementing a new type only requires adding a new mapper and calling this in init().
func RegisterMapper(mapper Mapper) {
	if err := globalRegistry.register(mapper); err != nil {
		panic(err)
	}
}

func (r *registry) register(mapper Mapper) error {
	if mapper == nil {
		return fmt.Errorf("document type mapper is nil")
	}

	name := normalizeName(mapper.Name())
	if name == "" {
		return fmt.Errorf("document type mapper name is empty")
	}

	dtoType := mapper.DTOType()
	modelType := mapper.ModelType()
	protoType := mapper.ProtoType()
	if dtoType == "" || modelType == "" || protoType == pb.DocumentType_DOCUMENT_TYPE_UNSPECIFIED {
		return fmt.Errorf("document type mapper %q has invalid mapping", name)
	}

	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("document type mapper %q already registered", name)
	}
	if _, exists := r.byProto[protoType]; exists {
		return fmt.Errorf("document type proto mapping %v already registered", protoType)
	}
	if _, exists := r.byDTO[dtoType]; exists {
		return fmt.Errorf("document type dto mapping %q already registered", dtoType)
	}
	if _, exists := r.byModel[modelType]; exists {
		return fmt.Errorf("document type model mapping %q already registered", modelType)
	}

	r.byName[name] = mapper
	r.byProto[protoType] = mapper
	r.byDTO[dtoType] = mapper
	r.byModel[modelType] = mapper
	return nil
}

func (r *registry) findByName(name string) Mapper {
	return r.byName[normalizeName(name)]
}

func (r *registry) findByProto(t pb.DocumentType) Mapper {
	return r.byProto[t]
}

func (r *registry) findByDTO(t dto.DocumentType) Mapper {
	if mapper, ok := r.byDTO[t]; ok {
		return mapper
	}
	return r.byName[normalizeName(string(t))]
}

func (r *registry) findByModel(t model.DocumentType) Mapper {
	if mapper, ok := r.byModel[t]; ok {
		return mapper
	}
	return r.byName[normalizeName(string(t))]
}

func (r *registry) defaultMapper() Mapper {
	if mapper := r.findByName(defaultTypeName); mapper != nil {
		return mapper
	}
	for _, mapper := range r.byName {
		return mapper
	}
	return nil
}

func normalizeName(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	switch normalized {
	case "slide":
		return "markdown_slide"
	case "yoreseerichtext", "yoresee-rich-text":
		return "yoresee_rich_text"
	default:
		return normalized
	}
}

func DefaultDTOType() dto.DocumentType {
	mapper := globalRegistry.defaultMapper()
	if mapper == nil {
		return dto.DocumentType(defaultTypeName)
	}
	return mapper.DTOType()
}

func IsSupportedDTOType(t dto.DocumentType) bool {
	return globalRegistry.findByDTO(t) != nil
}

func FromProtoType(t pb.DocumentType) dto.DocumentType {
	if mapper := globalRegistry.findByProto(t); mapper != nil {
		return mapper.DTOType()
	}
	return DefaultDTOType()
}

func ToProtoType(t dto.DocumentType) pb.DocumentType {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.ProtoType()
	}
	return pb.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
}

func ToModelType(t dto.DocumentType) model.DocumentType {
	if mapper := globalRegistry.findByDTO(t); mapper != nil {
		return mapper.ModelType()
	}
	return model.DocumentType(DefaultDTOType())
}

func FromModelType(t model.DocumentType) dto.DocumentType {
	if mapper := globalRegistry.findByModel(t); mapper != nil {
		return mapper.DTOType()
	}
	normalized := normalizeName(string(t))
	if normalized != "" {
		return dto.DocumentType(normalized)
	}
	return DefaultDTOType()
}
