package doc_type_mapper

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type markdownSlideAdapter struct{}

func (markdownSlideAdapter) Name() string {
	return string(dto.DocumentType_MarkdownSlide)
}

func (markdownSlideAdapter) ProtoType() pb.DocumentType {
	return pb.DocumentType_DOCUMENT_TYPE_MARKDOWN_SLIDE
}

func (markdownSlideAdapter) DTOType() dto.DocumentType {
	return dto.DocumentType_MarkdownSlide
}

func (markdownSlideAdapter) ModelType() model.DocumentType {
	return model.DocumentType_MarkdownSlide
}

func init() { RegisterMapper(markdownSlideAdapter{}) }
