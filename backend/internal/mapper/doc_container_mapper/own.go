package doc_container_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type ownAdapter struct{}

func (ownAdapter) Name() string {
	return string(dto.ContainerType_Own)
}

func (ownAdapter) ProtoType() pb.CreateDocumentContainerType {
	return pb.CreateDocumentContainerType_CREATE_DOCUMENT_CONTAINER_TYPE_OWN
}

func (ownAdapter) DTOType() dto.ContainerType {
	return dto.ContainerType_Own
}

func (ownAdapter) ModelType() model.ContainerType {
	return model.ContainerType_Own
}

func init() { RegisterMapper(ownAdapter{}) }
