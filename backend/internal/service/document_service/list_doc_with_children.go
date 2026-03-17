package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"

	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
)

func (s *DocumentService) listDocumentsWithChildren(req *internal_dto.DocumentsListReq) ([]*dto.DocumentMetaResponse, int64, error) {
	listOp, err := s.buildListDocumentsOperation(req)
	if err != nil {
		return nil, 0, err
	}
	models, total, err := listOp.ExecWithTotal()
	if err != nil {
		return nil, 0, err
	}

	if req.Options != nil && req.Options.IncludeChildren && req.Options.Recursive {
		parentIDs := make([]int64, len(models))
		for i := range models {
			parentIDs[i] = models[i].ID
		}

		allDescendants, err := s.getAllDescendantsByParentIDs(parentIDs)
		if err != nil {
			return nil, 0, err
		}

		docs := s.buildDocumentTree(models, allDescendants)
		return docs, total, nil
	}

	docs := make([]*dto.DocumentMetaResponse, 0, len(models))
	for i := range models {
		docs = append(docs, s.ConvertToDocumentResponse(&models[i]))
	}

	if req.Options != nil && req.Options.IncludeChildren {
		for i := range models {
			childDocs, err := s.getChildDocuments(models[i].ID, req.Options)
			if err == nil {
				docs[i].Children = childDocs
				if len(childDocs) > 0 {
					docs[i].HasChildren = true
				}
			}
		}
	}

	return docs, total, nil
}

func (s *DocumentService) ListDocumentsWithChildrenByExternal(ctx context.Context, req *dto.ListDocumentsByExternalReq) ([]*dto.DocumentMetaResponse, int64, error) {
	var userID *int64
	if req == nil {
		return nil, 0, status.StatusInternalParamsError
	}
	if req.ExternalArgs != nil && req.ExternalArgs.UserExternalID != nil && *req.ExternalArgs.UserExternalID != "" {
		id, err := s.userRepo.GetIDByExternalID(*req.ExternalArgs.UserExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		userID = &id
	}

	var parentID int64
	if req.ExternalArgs != nil && req.ExternalArgs.RootDocumentExternalID != nil && *req.ExternalArgs.RootDocumentExternalID != "" {
		doc, err := s.documentRepo.GetByExternalID(*req.ExternalArgs.RootDocumentExternalID).Exec(ctx)
		if err != nil {
			return nil, 0, err
		}
		parentID = doc.ID
	}

	var knowledgeID *int64
	if req.ExternalArgs != nil && req.ExternalArgs.KnowledgeExternalID != nil && *req.ExternalArgs.KnowledgeExternalID != "" {
		id, err := s.kbRepo.GetIDByExternalID(*req.ExternalArgs.KnowledgeExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		knowledgeID = &id
	}

	var listOwnDoc bool
	if req.ExternalArgs != nil {
		listOwnDoc = req.ExternalArgs.ListOwnDoc
	}

	return s.listDocumentsWithChildren(
		&internal_dto.DocumentsListReq{
			MetaArgs: &internal_dto.DocumentsListMetaArgs{
				UserID:      userID,
				ParentID:    &parentID, // default to root
				KnowledgeID: knowledgeID,
				ListOwnDoc:  listOwnDoc,
			},
			FilterArgs: req.FilterArgs,
			SortArgs:   req.SortArgs,
			Pagination: req.Pagination,
			Options:    req.Options,
		},
	)
}

func (s *DocumentService) getChildDocuments(parentID int64, options *dto.RecursiveOptions) ([]*dto.DocumentMetaResponse, error) {
	models, _, err := s.documentRepo.ListDocuments(&model.Document{}).
		WithParentID(&parentID).
		WithSort("created_at", false).
		ExecWithTotal()
	if err != nil {
		return nil, err
	}

	childResponses := make([]*dto.DocumentMetaResponse, len(models))
	for i, model := range models {
		childResponses[i] = s.ConvertToDocumentResponse(&model)

		if options != nil && options.Recursive && (options.Depth <= 0 || options.Depth > 1) {
			subOptions := &dto.RecursiveOptions{
				IncludeChildren: options.IncludeChildren,
				Recursive:       options.Recursive,
				Depth:           options.Depth - 1,
			}
			grandChildren, err := s.getChildDocuments(model.ID, subOptions)
			if err == nil {
				childResponses[i].Children = grandChildren
				if len(grandChildren) > 0 {
					childResponses[i].HasChildren = true
				}
			}
		}
	}

	return childResponses, nil
}

func (s *DocumentService) getAllDescendantsByParentIDs(parentIDs []int64) ([]model.Document, error) {
	if len(parentIDs) == 0 {
		return []model.Document{}, nil
	}

	allDocs := make([]model.Document, 0)
	seen := make(map[int64]bool)

	for _, rootParentID := range parentIDs {
		docs, err := s.documentRepo.GetSubtree(rootParentID).Exec()
		if err != nil {
			return nil, err
		}

		for _, doc := range docs {
			if !seen[doc.ID] {
				seen[doc.ID] = true
				allDocs = append(allDocs, doc)
			}
		}
	}

	return allDocs, nil
}
