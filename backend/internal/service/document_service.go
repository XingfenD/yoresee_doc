package service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/bytedance/gg/gslice"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DocumentService struct {
	documentRepo   *repository.DocumentRepository
	userRepo       *repository.UserRepository
	kbRepo         *repository.KnowledgeBaseRepository
	docVersionRepo *repository.DocumentVersionRepository
	snapshotRepo   *repository.DocumentYjsSnapshotRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: &repository.DocumentRepo,
		userRepo:     repository.UserRepo,
		kbRepo:       repository.KnowledgeBaseRepo,
		snapshotRepo: repository.DocumentYjsSnapshotRepo,
	}
}

var DocumentSvc = NewDocumentService()

func (s *DocumentService) ConvertToDocumentResponse(doc *model.Document) *dto.DocumentMetaResponse {
	return dto.NewDocumentMetaResponseFromModel(doc)
}

func (s *DocumentService) GetDocumentByExternalID(ctx context.Context, externalID string) (*dto.DocumentResponse, error) {
	docModel, err := s.documentRepo.GetByExternalID(externalID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return dto.NewDocumentResponseFromModel(docModel), nil
}

func (s *DocumentService) GetDocumentYjsSnapshot(ctx context.Context, docExternalID string) ([]byte, error) {
	docID, err := repository.DocumentRepo.GetIDByExternalID(docExternalID).Exec(ctx)
	if err != nil {
		return nil, status.StatusDocumentNotFound
	}

	snapshot, err := s.snapshotRepo.GetByDocID(docID).Exec()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return snapshot.YjsState, nil
}

func (s *DocumentService) SaveDocumentYjsSnapshot(ctx context.Context, docExternalID string, state []byte) error {
	if len(state) == 0 {
		return status.StatusParamError
	}
	docID, err := repository.DocumentRepo.GetIDByExternalID(docExternalID).Exec(ctx)
	if err != nil {
		return status.StatusDocumentNotFound
	}

	return s.snapshotRepo.Save(docID, state).Exec()
}

func (s *DocumentService) buildListDocumentsOperation(req *documentsListReq) (*repository.DocumentsListOperation, error) {
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
			WithListOwnDoc(req.MetaArgs.ListOwnDoc)
	}
	if req.FilterArgs != nil {
		listOp = listOp.WithTitleKeyword(req.FilterArgs.TitleKeyword).
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

func (s *DocumentService) buildDocumentTree(rootDocs []model.Document, allDescendants []model.Document) []*dto.DocumentMetaResponse {
	docMap := make(map[int64]*dto.DocumentMetaResponse)
	var rootResponses []*dto.DocumentMetaResponse

	for i := range rootDocs {
		resp := s.ConvertToDocumentResponse(&rootDocs[i])
		docMap[rootDocs[i].ID] = resp
		rootResponses = append(rootResponses, resp)
	}

	for i := range allDescendants {
		doc := &allDescendants[i]
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

func (s *DocumentService) listDocumentsWithChildren(req *documentsListReq) ([]*dto.DocumentMetaResponse, int64, error) {
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
		id, err := repository.UserRepo.GetIDByExternalID(*req.ExternalArgs.UserExternalID).Exec()
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
		id, err := repository.KnowledgeBaseRepo.GetIDByExternalID(*req.ExternalArgs.KnowledgeExternalID).Exec()
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
		&documentsListReq{
			MetaArgs: &documentsListMetaArgs{
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

func (s *DocumentService) Create(ctx context.Context, req *dto.CreateDocumentReq) (*dto.CreateDocumentResponse, error) {
	if err := validateCreateDocumentReq(req); err != nil {
		return nil, err
	}

	// TODO: redis support
	docExternalID := utils.GenerateExternalID(utils.ExternalIDContextDocument)
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		docModel := &model.Document{
			ExternalID: docExternalID,
			Title:      req.Title,
			Type:       model.DocumentType(req.Type),
			Summary:    "",
			Content:    "",
			Path:       "0",
			Depth:      0,
		}

		// query user_id
		userID, err := s.userRepo.GetIDByExternalID(*req.CreatorExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}

		docModel.UserID = userID

		// query knowledge_id
		if !req.CreateAsOwnDoc {
			kbID, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeExternalID).WithTx(tx).Exec()
			if err != nil {
				return status.StatusKnowledgeBaseNotFound
			}
			docModel.KnowledgeID = &kbID
			docModel.ContainerType = model.ContainerType_KnowledgeBase
		} else {
			docModel.ContainerType = model.ContainerType_Own
		}

		// TODO: permission check for knowledgebase

		// query parent_id
		if req.ParentExternalID != nil {
			parentDocID, err := s.documentRepo.GetIDByExternalID(*req.ParentExternalID).WithTx(tx).Exec(ctx)
			if err != nil {
				return status.StatusDocumentNotFound
			}
			docModel.ParentID = parentDocID
		}

		// create document meta
		err = s.documentRepo.Create(docModel).WithTx(tx).Exec()
		if err != nil {
			status.GenErrWithCustomMsg(status.StatusWriteDBError, "create document meta failed")
			return status.StatusWriteDBError
		}
		if err := s.documentRepo.UpdatePathDepth(docModel.ID, docModel.ParentID).WithTx(tx).Exec(); err != nil {
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "update document path depth failed")
		}

		// create doc version
		ver := &model.DocumentVersion{
			DocumentID:    docModel.ID,
			Content:       docModel.Content,
			UserID:        userID,
			Title:         docModel.Title,
			ChangeSummary: "Create the document",
		}
		err = s.docVersionRepo.Create(ver).WithTx(tx).Exec()
		if err != nil {
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create version failed")
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("[Service layer: DocumentService] Create err: %+v", err)
		return nil, status.StatusWriteDBError
	}
	return &dto.CreateDocumentResponse{
		ExternalID: docExternalID,
	}, nil
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

func (s *DocumentService) Update(ctx context.Context, req *dto.UpdateDocumentRequest) (bool, error) {
	if err := validateUpdateDocumentReq(req); err != nil {
		return false, err
	}

	err := utils.WithTransaction(func(tx *gorm.DB) error {
		oldDoc, err := s.documentRepo.GetByExternalID(req.ExternalID).WithTx(tx).Exec(ctx)
		if err != nil {
			return status.StatusDocumentNotFound
		}

		docModel := &model.Document{
			ID: oldDoc.ID,
		}
		op := s.documentRepo.Update(docModel).WithTx(tx)
		if req.Content != nil {
			docModel.Content = *req.Content
			op.UpdateContent()
		}
		if req.Title != nil {
			docModel.Title = *req.Title
		}

		moved := false
		var newParentID int64
		if req.ParentExternalID != nil {
			parentID, err := s.documentRepo.GetIDByExternalID(*req.ParentExternalID).Exec(ctx)
			if err != nil {
				return status.StatusDocumentNotFound
			}
			docModel.ParentID = parentID
			newParentID = parentID
			moved = true
			op = op.UpdateParentID()
		}

		// create version
		versionModel := &model.DocumentVersion{
			DocumentID:    oldDoc.ID,
			Title:         oldDoc.Title,
			Content:       oldDoc.Content,
			UserID:        oldDoc.UserID,
			ChangeSummary: "",
		}

		if err := s.docVersionRepo.Create(versionModel).WithTx(tx).Exec(); err != nil {
			return err
		}

		// update kb relation
		if req.KnowledgeBaseExternalID != nil && !req.MoveAsOwn {
			kbID, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeBaseExternalID).Exec()
			if err != nil {
				return status.StatusKnowledgeBaseNotFound
			}
			docModel.KnowledgeID = &kbID
			op = op.UpdateKnowledgeID()
		}

		if req.MoveAsOwn {
			docModel.KnowledgeID = nil
			op = op.UpdateKnowledgeID()
		}

		if err := op.Exec(); err != nil {
			return err
		}
		if moved {
			if err := s.documentRepo.MoveSubtree(oldDoc.ID, newParentID).WithTx(tx).Exec(); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return false, err
	}

	// TODO: redis support

	return true, nil
}
