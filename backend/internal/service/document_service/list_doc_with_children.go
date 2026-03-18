package document_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/status"

	internal_dto "github.com/XingfenD/yoresee_doc/internal/service/dto"
)

func (s *DocumentService) listDocuments(ctx context.Context, req *internal_dto.DocumentsListReq) ([]*dto.DocumentMetaResponse, int64, error) {
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

	getOp := s.getDocumentsWithDescendants(parentIDs).
		WithDirectoryOnly(req.DirectoryOnly)
	if req.Options != nil {
		getOp = getOp.WithDepth(req.Options.Depth)
	}

	allDescendants, err := getOp.Exec(ctx)
	if err != nil {
		return nil, 0, err
	}

	docs := s.buildDocumentTree(models, allDescendants)
	return docs, total, nil
}

func (s *DocumentService) ListDocumentsByExternal(ctx context.Context, req *dto.ListDocumentsByExternalReq) ([]*dto.DocumentMetaResponse, int64, error) {
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

	return s.listDocuments(
		ctx,
		&internal_dto.DocumentsListReq{
			MetaArgs: &internal_dto.DocumentsListMetaArgs{
				UserID:      userID,
				ParentID:    &parentID, // default to root
				KnowledgeID: knowledgeID,
			},
			ListDocumentsBaseArgs: dto.ListDocumentsBaseArgs{
				ListOwnDoc:    req.ListOwnDoc,
				DirectoryOnly: req.DirectoryOnly,
			},
			FilterArgs: req.FilterArgs,
			SortArgs:   req.SortArgs,
			Pagination: req.Pagination,
			Options:    req.Options,
		},
	)
}

type GetDocumentsWithDescendantsOption struct {
	DirectoryOnly bool
}

type GetDocumentsWithDescendantsOperation struct {
	svc       *DocumentService
	parentIDs []int64
	depth     *int

	directoryOnly bool
}

func (s *DocumentService) getDocumentsWithDescendants(parentIDs []int64) *GetDocumentsWithDescendantsOperation {
	return &GetDocumentsWithDescendantsOperation{
		svc:       s,
		parentIDs: parentIDs,
	}
}

func (op *GetDocumentsWithDescendantsOperation) WithDepth(depth *int) *GetDocumentsWithDescendantsOperation {
	op.depth = depth
	return op
}

func (op *GetDocumentsWithDescendantsOperation) WithDirectoryOnly(with bool) *GetDocumentsWithDescendantsOperation {
	op.directoryOnly = with
	return op
}

func (op *GetDocumentsWithDescendantsOperation) Exec(ctx context.Context) ([]*model.Document, error) {
	if len(op.parentIDs) == 0 {
		return []*model.Document{}, nil
	}

	allDocs := make([]*model.Document, 0)
	seen := make(map[int64]bool)

	for _, rootParentID := range op.parentIDs {
		docs, err := op.svc.documentRepo.GetSubtree(rootParentID).
			WithDepth(op.depth).
			WithDirectoryOnly(op.directoryOnly).
			Exec(ctx)
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
