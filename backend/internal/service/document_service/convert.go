package document_service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository/document_repo"
	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

func (s *DocumentService) ConvertToDocumentResponse(doc *model.Document) *dto.DocumentMetaResponse {
	return dto.NewDocumentMetaResponseFromModel(doc)
}

func (s *DocumentService) buildListDocumentsOperation(req *internal_dto.DocumentsListReq) (*document_repo.DocumentsListOperation, error) {
	if s == nil || s.documentRepo == nil {
		return nil, status.StatusServiceInternalError
	}
	if req == nil {
		return nil, status.StatusInternalParamsError
	}
	listOp := s.documentRepo.ListDocuments(&model.Document{})
	if req.MetaArgs != nil {
		listOp = listOp.WithUserID(req.MetaArgs.UserID).
			WithParentID(req.MetaArgs.ParentID).
			WithKnowledgeID(req.MetaArgs.KnowledgeID).
			WithListOwnDoc(req.ListOwnDoc).
			WithDirectoryOnly(req.DirectoryOnly)
	}
	if req.SearchDocIDs != nil {
		listOp = listOp.WithIDs(req.SearchDocIDs)
	}
	if req.FilterArgs != nil {
		titleKeyword := req.FilterArgs.TitleKeyword
		if req.SearchDocIDs != nil {
			titleKeyword = nil
		}
		listOp = listOp.WithTitleKeyword(titleKeyword).
			WithType(req.FilterArgs.DocType).
			WithStatus(req.FilterArgs.Status).
			WithTags(req.FilterArgs.Tags).
			WithCreateTimeRange(req.FilterArgs.CreateTimeRangeStart, req.FilterArgs.CreateTimeRangeEnd).
			WithUpdateTimeRange(req.FilterArgs.UpdateTimeRangeStart, req.FilterArgs.UpdateTimeRangeEnd)
	}
	listOp = listOp.WithSort(req.SortArgs.Field, req.SortArgs.Desc).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize)

	return listOp, nil
}

func (s *DocumentService) buildDocumentTree(rootDocs []model.Document, allDescendants []*model.Document) []*dto.DocumentMetaResponse {
	docMap := make(map[int64]*dto.DocumentMetaResponse)
	var rootResponses []*dto.DocumentMetaResponse

	for i := range rootDocs {
		resp := s.ConvertToDocumentResponse(&rootDocs[i])
		docMap[rootDocs[i].ID] = resp
		rootResponses = append(rootResponses, resp)
	}

	for _, doc := range allDescendants {
		childResp := s.ConvertToDocumentResponse(doc)
		docMap[doc.ID] = childResp

		if doc.ParentID != 0 {
			parentResp, exists := docMap[doc.ParentID]
			if exists {
				parentResp.Children = append(parentResp.Children, childResp)
				parentResp.HasChildren = true
			}
		}
	}

	return rootResponses
}
