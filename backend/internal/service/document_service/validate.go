package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	doc_container_mapper "github.com/XingfenD/yoresee_doc/internal/mapper/doc_container_mapper"
	"github.com/XingfenD/yoresee_doc/internal/mapper/doc_type_mapper"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func validateCreateDocumentReq(req *dto.CreateDocumentReq) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if !doc_container_mapper.IsSupportedDTOType(req.ContainerType) {
		return status.GenErrWithCustomMsg(status.StatusParamError, "invalid document container type")
	}
	if req.ContainerType == dto.ContainerType_Own && req.KnowledgeExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID not nil when container_type is own")
	}
	if req.ContainerType == dto.ContainerType_KnowledgeBase && req.KnowledgeExternalID == nil {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "KnowledgeExternalID is nil when container_type is knowledge_base")
	}

	if !doc_type_mapper.IsSupportedDTOType(req.Type) {
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

	if req.Content == nil && req.KnowledgeBaseExternalID == nil && req.ParentExternalID == nil && req.Title == nil && req.MoveToContainer == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update field")
	}

	if req.MoveToContainer != nil {
		if !doc_container_mapper.IsSupportedDTOType(*req.MoveToContainer) {
			return status.GenErrWithCustomMsg(status.StatusParamError, "invalid move_to_container")
		}
		if *req.MoveToContainer == dto.ContainerType_Own && req.KnowledgeBaseExternalID != nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalID not nil when moving to own")
		}
		if *req.MoveToContainer == dto.ContainerType_KnowledgeBase && req.KnowledgeBaseExternalID == nil {
			return status.GenErrWithCustomMsg(status.StatusParamError, "KnowledgeBaseExternalID is nil when moving to knowledge_base")
		}
	}
	if req.MoveToContainer == nil && req.KnowledgeBaseExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "MoveToContainer is nil when KnowledgeBaseExternalID provided")
	}

	if req.MoveToContainer != nil && *req.MoveToContainer == dto.ContainerType_Own && req.ParentExternalID != nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "ParentExternalID not nil when moving to own")
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

func validateUpdateTemplateSettingsReq(req *dto.UpdateTemplateSettingsRequest) error {
	if req == nil {
		return status.StatusInternalParamsError
	}
	if req.UserExternalID == "" {
		return status.GenErrWithCustomMsg(status.StatusParamError, "user external id is empty")
	}
	if req.TemplateID <= 0 {
		return status.GenErrWithCustomMsg(status.StatusParamError, "template id is invalid")
	}
	if req.Name == nil && req.Description == nil && req.IsPublic == nil {
		return status.GenErrWithCustomMsg(status.StatusParamError, "no update field")
	}
	return nil
}
