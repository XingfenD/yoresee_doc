package doc_type_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type tableAdapter struct{}

func (tableAdapter) Name() string {
	return string(dto.DocumentType_Table)
}

func (tableAdapter) ProtoType() pb.DocumentType {
	return pb.DocumentType_DOCUMENT_TYPE_TABLE
}

func (tableAdapter) DTOType() dto.DocumentType {
	return dto.DocumentType_Table
}

func (tableAdapter) ModelType() model.DocumentType {
	return model.DocumentType_Table
}

func init() { RegisterMapper(tableAdapter{}) }
