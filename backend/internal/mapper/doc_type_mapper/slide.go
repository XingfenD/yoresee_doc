package doc_type_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type slideAdapter struct{}

func (slideAdapter) Name() string {
	return string(dto.DocumentType_Slide)
}

func (slideAdapter) ProtoType() pb.DocumentType {
	return pb.DocumentType_DOCUMENT_TYPE_SLIDE
}

func (slideAdapter) DTOType() dto.DocumentType {
	return dto.DocumentType_Slide
}

func (slideAdapter) ModelType() model.DocumentType {
	return model.DocumentType_Slide
}

func init() { RegisterMapper(slideAdapter{}) }
