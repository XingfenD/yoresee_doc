package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"

	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
)

func (s *DocumentService) listDocuments(req *internal_dto.DocumentsListReq) ([]*dto.DocumentMetaResponse, int64, error) {
	listOp, err := s.buildListDocumentsOperation(req)
	if err != nil {
		return nil, 0, err
	}
	models, total, err := listOp.ExecWithTotal()
	if err != nil {
		return nil, 0, err
	}

	// if children is not needed or depth less than 1, return directly
	if req.Options == nil || !req.Options.IncludeChildren || (req.Options.Depth != nil && *req.Options.Depth < 1) {
		return s.buildDocumentTree(models, nil), total, nil
	}

	parentIDs := make([]int64, len(models))
	for i := range models {
		parentIDs[i] = models[i].ID
	}

	allDescendants, err := s.getDocumentsWithDescendants(parentIDs, req.Options.Depth)
	if err != nil {
		return nil, 0, err
	}

	docs := s.buildDocumentTree(models, allDescendants)
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

	return s.listDocuments(
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

func (s *DocumentService) getDocumentsWithDescendants(parentIDs []int64, depth *int) ([]model.Document, error) {
	if len(parentIDs) == 0 {
		return []model.Document{}, nil
	}

	allDocs := make([]model.Document, 0)
	seen := make(map[int64]bool)

	for _, rootParentID := range parentIDs {
		docs, err := s.documentRepo.GetSubtree(rootParentID).WithDepth(depth).Exec()
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
