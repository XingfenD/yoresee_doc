package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/bytedance/gg/gslice"
)

func validateCreateDocumentReq(req *dto.CreateDocumentReq) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.CreateAsOwnDoc && req.KnowledgeExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID not nil when CreateAsOwnDoc")
	}
	if !req.CreateAsOwnDoc && req.KnowledgeExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID is nil when not CreateAsOwnDoc")
	}

	availableTypes := []dto.DocumentType{dto.DocumentType_Markdown}
	if !gslice.Contains(availableTypes, req.Type) {
		return status.GenErrWithCustomMsg(status.StatusParamError, "invalid document type")
	}

	return nil
}

func validateUpdateDocumentReq(req *dto.UpdateDocumentRequest) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.ExternalID == "" {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "external_id is zero value")
	}

	if req.Content == nil && req.KnowledgeBaseExternalID == nil && req.ParentExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update field")
	}

	if req.MoveAsOwn && req.ParentExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "ParentExternalID not nil when moving as own")
	}

	return nil
}

func validateUpdateDocumentMetaReq(req *dto.UpdateDocumentMetaRequest) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.ExternalID == "" {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "external_id is zero value")
	}
	if req.Title == nil && req.Summary == nil && req.Tags == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update field")
	}
	return nil
}

func validateCreateTemplateReq(req *dto.CreateTemplateRequest) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.UserExternalID == "" {
		return status.GenErrWithCustomMsg(status.StatusParamError, "user external id is empty")
	}
	if req.TemplateContent == "" {
		return status.GenErrWithCustomMsg(status.StatusParamError, "template content is empty")
	}
	switch req.TargetContainer {
	case dto.TemplateContainerOwn, dto.TemplateContainerPublic:
		return nil
	case dto.TemplateContainerKnowledgeBase:
		if req.KnowledgeBaseExternalID == nil || *req.KnowledgeBaseExternalID == "" {
			return status.GenErrWithCustomMsg(status.StatusParamError, "knowledge_base_id is empty")
		}
		return nil
	default:
		return status.GenErrWithCustomMsg(status.StatusParamError, "invalid template container")
	}
}
