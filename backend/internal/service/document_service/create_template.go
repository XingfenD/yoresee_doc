package document_service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/mapper/template_container_mapper"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
)

type templatePayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Tags        []string `json:"tags"`
}

func (s *DocumentService) CreateTemplate(ctx context.Context, req *dto.CreateTemplateRequest) error {
	if err := validateCreateTemplateReq(req); err != nil {
		logrus.Errorf("[Service layer: DocumentService] validateCreateTemplateReq failed, err=%+v", err)
		return status.GenErrWithCustomMsg(err, "invalid create template request")
	}
	_ = ctx

	userID, err := s.userRepo.GetIDByExternalID(req.UserExternalID).Exec()
	if err != nil {
		return status.StatusUserNotFound
	}

	var kbID *int64
	if template_container_mapper.RequiresKnowledgeBaseID(req.TargetContainer) {
		id, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeBaseExternalID).Exec()
		if err != nil {
			return status.StatusKnowledgeBaseNotFound
		}
		kbID = &id
	}

	name, description, content, tags := parseTemplateContent(req.TemplateContent)
	if strings.TrimSpace(content) == "" {
		return status.GenErrWithCustomMsg(status.StatusParamError, "template content is empty")
	}

	templateType := req.Type
	if !doc_type_mapper.IsSupportedDTOType(templateType) {
		templateType = doc_type_mapper.DefaultDTOType()
	}

	scope := template_container_mapper.ToScope(req.TargetContainer)
	isPublic := template_container_mapper.ToIsPublic(req.TargetContainer)
	templateModel := &model.Template{
		Name:            name,
		Description:     description,
		DocumentType:    doc_type_mapper.ToModelType(templateType),
		Content:         content,
		UserID:          userID,
		Scope:           scope,
		KnowledgeBaseID: kbID,
		IsPublic:        isPublic,
		Tags:            tags,
	}

	if err := s.templateRepo.Create(templateModel).Exec(); err != nil {
		return status.StatusWriteDBError
	}

	return nil
}

func parseTemplateContent(raw string) (name, description, content string, tags []string) {
	content = raw
	var payload templatePayload
	if err := json.Unmarshal([]byte(raw), &payload); err == nil {
		if payload.Content != "" {
			content = payload.Content
		}
		name = payload.Name
		description = payload.Description
		tags = payload.Tags
	}
	if strings.TrimSpace(name) == "" {
		name = deriveTemplateName(content)
	}
	if len(name) > 100 {
		name = truncateRunes(name, 100)
	}
	return name, description, content, tags
}

func deriveTemplateName(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		trimmed = strings.TrimLeft(trimmed, "#*-+ ")
		trimmed = strings.TrimSpace(trimmed)
		if trimmed == "" {
			continue
		}
		return trimmed
	}
	return "Untitled Template"
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
