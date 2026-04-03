package doc_type_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type markdownAdapter struct{}

func (markdownAdapter) Name() string {
	return string(dto.DocumentType_Markdown)
}

func (markdownAdapter) ProtoType() pb.DocumentType {
	return pb.DocumentType_DOCUMENT_TYPE_MARKDOWN
}

func (markdownAdapter) DTOType() dto.DocumentType {
	return dto.DocumentType_Markdown
}

func (markdownAdapter) ModelType() model.DocumentType {
	return model.DocumentType_Markdown
}

func init() { RegisterMapper(markdownAdapter{}) }
