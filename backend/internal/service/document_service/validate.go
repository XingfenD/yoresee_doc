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
