package doc_type_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type yoreseeRichTextAdapter struct{}

func (yoreseeRichTextAdapter) Name() string {
	return string(dto.DocumentType_YoreseeRichText)
}

func (yoreseeRichTextAdapter) ProtoType() pb.DocumentType {
	return pb.DocumentType_DOCUMENT_TYPE_YORESEE_RICH_TEXT
}

func (yoreseeRichTextAdapter) DTOType() dto.DocumentType {
	return dto.DocumentType_YoreseeRichText
}

func (yoreseeRichTextAdapter) ModelType() model.DocumentType {
	return model.DocumentType_YoreseeRichText
}

func init() { RegisterMapper(yoreseeRichTextAdapter{}) }
