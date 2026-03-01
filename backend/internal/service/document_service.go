package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
)

type DocumentService struct {
	documentRepo *repository.DocumentRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: repository.DocumentRepo,
	}
}

var DocumentSvc = NewDocumentService()

func (s *DocumentService) ConvertToDocumentResponse(doc *model.DocumentMeta) *dto.DocumentResponse {
	return dto.NewDocumentResponseFromModel(doc)
}

func (s *DocumentService) GetDocumentByExternalID(externalID string) (*model.DocumentMeta, error) {
	return s.documentRepo.GetByExternalID(externalID).Exec()
}

func (s *DocumentService) GetDocumentContent(documentID int64) (string, error) {
	return s.documentRepo.GetContent(documentID).Exec()
}

// func (s *DocumentService) CheckDocumentPermission(userID int64, documentID int64, requiredPermission string) (bool, error) {
// 	cacheKey := fmt.Sprintf("permission:user:%d:doc:%d:%s", userID, documentID, requiredPermission)
// 	ctx := context.Background()

// 	cachedResult, err := storage.GetCache(ctx, cacheKey)
// 	if err == nil {
// 		return cachedResult == "true", nil
// 	}

// 	permissionCheck := &dto.PermissionCheck{
// 		Resource: dto.Resource{
// 			Type: string(model.ResourceTypeDocument),
// 			ID:   fmt.Sprintf("%d", documentID),
// 		},
// 		Permission: requiredPermission,
// 	}
// 	allowed, err := s.permissionSvc.CheckPermission(userID, permissionCheck)
// 	if err != nil {
// 		return false, err
// 	}

// 	err = storage.SetCache(ctx, cacheKey, fmt.Sprintf("%v", allowed), time.Hour)
// 	if err != nil {
// 		logrus.Info("Set permission cache failed: %v\n", err)
// 	}

// 	return allowed, nil
// }

func (s *DocumentService) GetDocumentWithContent(docExternalID string) (*model.DocumentMeta, string, error) {
	document, err := s.GetDocumentByExternalID(docExternalID)
	if err != nil {
		return nil, "", err
	}

	content, err := s.GetDocumentContent(document.ID)
	if err != nil {
		return document, "", err
	}

	return document, content, nil
}

type DocumentsListReq struct {
	MetaArgs   *DocumentsListMetaArgs   `json:"meta_args"`
	FilterArgs *DocumentsListFilterArgs `json:"filter_args"`
	SortArgs   SortArgs                 `json:"sort_args"`
	Pagination Pagination               `json:"pagination"`
	Options    *ListDocumentsOptions    `json:"options"`
}

type DocumentsListMetaArgs struct {
	UserID      *int64 `json:"user_id"`
	ParentID    *int64 `json:"parent_id"`
	KnowledgeID *int64 `json:"knowledge_id"`
}

type DocumentsListExternalArgs struct {
	UserExternalID         *string `json:"user_external_id"`
	RootDocumentExternalID *string `json:"root_document_external_id"`
	KnowledgeExternalID    *string `json:"knowledge_external_id"`
}

type DocumentsListFilterArgs struct {
	TitleKeyword         *string  `json:"title_keyword"`
	DocType              *string  `json:"doc_type"`
	Status               *int     `json:"status"`
	Tags                 []string `json:"tags"`
	CreateTimeRangeStart *string  `json:"create_time_range_start"`
	CreateTimeRangeEnd   *string  `json:"create_time_range_end"`
	UpdateTimeRangeStart *string  `json:"update_time_range_start"`
	UpdateTimeRangeEnd   *string  `json:"update_time_range_end"`
}

func (s *DocumentService) ListDocuments(req *DocumentsListReq) ([]*dto.DocumentResponse, int64, error) {
	models, total, err := s.documentRepo.ListDocuments(&model.DocumentMeta{}).
		WithUserID(req.MetaArgs.UserID).
		WithParentID(req.MetaArgs.ParentID).
		WithKnowledgeID(req.MetaArgs.KnowledgeID).
		WithTitleKeyword(req.FilterArgs.TitleKeyword).
		WithType(req.FilterArgs.DocType).
		WithStatus(req.FilterArgs.Status).
		WithTags(req.FilterArgs.Tags).
		WithCreateTimeRange(req.FilterArgs.CreateTimeRangeStart, req.FilterArgs.CreateTimeRangeEnd).
		WithUpdateTimeRange(req.FilterArgs.UpdateTimeRangeStart, req.FilterArgs.UpdateTimeRangeEnd).
		WithSort(req.SortArgs.Field, req.SortArgs.Desc).
		WithPagination(req.Pagination.Page, req.Pagination.PageSize).
		ExecWithTotal()

	if err != nil {
		return nil, 0, err
	}
	responses := make([]*dto.DocumentResponse, 0, len(models))
	for _, model := range models {
		responses = append(responses, s.ConvertToDocumentResponse(&model))
	}
	return responses, total, nil
}

type ListDocumentsOptions struct {
	IncludeChildren bool `json:"include_children"`
	Recursive       bool `json:"recursive"`
	Depth           int  `json:"depth"`
}

func (s *DocumentService) GetChildDocuments(parentID int64, options *ListDocumentsOptions) ([]*dto.DocumentResponse, error) {
	models, _, err := s.documentRepo.ListDocuments(&model.DocumentMeta{}).
		WithParentID(&parentID).
		WithSort("created_at", false).
		ExecWithTotal()
	if err != nil {
		return nil, err
	}

	childResponses := make([]*dto.DocumentResponse, len(models))
	for i, model := range models {
		childResponses[i] = s.ConvertToDocumentResponse(&model)

		if options != nil && options.Recursive && (options.Depth <= 0 || options.Depth > 1) {
			subOptions := &ListDocumentsOptions{
				IncludeChildren: options.IncludeChildren,
				Recursive:       options.Recursive,
				Depth:           options.Depth - 1,
			}
			grandChildren, err := s.GetChildDocuments(model.ID, subOptions)
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

func (s *DocumentService) ListDocumentsWithChildren(req *DocumentsListReq) ([]*dto.DocumentResponse, int64, error) {
	docs, total, err := s.ListDocuments(req)
	if err != nil {
		return nil, 0, err
	}
	if req.Options != nil && req.Options.IncludeChildren {
		models, _, err := s.documentRepo.ListDocuments(&model.DocumentMeta{}).
			WithUserID(req.MetaArgs.UserID).
			WithParentID(req.MetaArgs.ParentID).
			WithTitleKeyword(req.FilterArgs.TitleKeyword).
			WithType(req.FilterArgs.DocType).
			WithStatus(req.FilterArgs.Status).
			WithTags(req.FilterArgs.Tags).
			WithCreateTimeRange(req.FilterArgs.CreateTimeRangeStart, req.FilterArgs.CreateTimeRangeEnd).
			WithUpdateTimeRange(req.FilterArgs.UpdateTimeRangeStart, req.FilterArgs.UpdateTimeRangeEnd).
			WithSort(req.SortArgs.Field, req.SortArgs.Desc).
			WithPagination(req.Pagination.Page, req.Pagination.PageSize).
			ExecWithTotal()
		if err != nil {
			return nil, 0, err
		}

		for i := range docs {
			childDocs, err := s.GetChildDocuments(models[i].ID, req.Options)
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

type ListDocumentsByExternalReq struct {
	ExternalArgs *DocumentsListExternalArgs `json:"external_args"`
	FilterArgs   *DocumentsListFilterArgs   `json:"filter_args"`
	SortArgs     SortArgs                   `json:"sort_args"`
	Pagination   Pagination                 `json:"pagination"`
	Options      *ListDocumentsOptions      `json:"options"`
}

func (s *DocumentService) ListDocumentsWithChildrenByExternal(req *ListDocumentsByExternalReq) ([]*dto.DocumentResponse, int64, error) {
	var userID *int64
	if req.ExternalArgs.UserExternalID != nil && *req.ExternalArgs.UserExternalID != "" {
		id, err := repository.UserRepo.GetIDByExternalID(*req.ExternalArgs.UserExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		userID = &id
	}

	var parentID *int64
	if req.ExternalArgs.RootDocumentExternalID != nil && *req.ExternalArgs.RootDocumentExternalID != "" {
		doc, err := s.GetDocumentByExternalID(*req.ExternalArgs.RootDocumentExternalID)
		if err != nil {
			return nil, 0, err
		}
		parentID = &doc.ID
	}

	var knowledgeID *int64
	if req.ExternalArgs.KnowledgeExternalID != nil && *req.ExternalArgs.KnowledgeExternalID != "" {
		id, err := repository.KnowledgeBaseRepo.GetIDByExternalID(*req.ExternalArgs.KnowledgeExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		knowledgeID = &id
	}

	return s.ListDocumentsWithChildren(
		&DocumentsListReq{
			MetaArgs: &DocumentsListMetaArgs{
				UserID:      userID,
				ParentID:    parentID,
				KnowledgeID: knowledgeID,
			},
			FilterArgs: req.FilterArgs,
			SortArgs:   req.SortArgs,
			Pagination: req.Pagination,
			Options:    req.Options,
		},
	)
}
