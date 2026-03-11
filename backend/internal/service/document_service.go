package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DocumentService struct {
	documentRepo      *repository.DocumentRepository
	userRepo          *repository.UserRepository
	kbRepo            *repository.KnowledgeBaseRepository
	docKBRelationRepo *repository.DocKnowledgeRelationRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: repository.DocumentRepo,
		userRepo:     repository.UserRepo,
		kbRepo:       repository.KnowledgeBaseRepo,
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

func (s *DocumentService) buildListDocumentsOperation(req *documentsListReq) (*repository.DocumentsListOperation, error) {
	if s == nil || s.documentRepo == nil {
		return nil, status.StatusServiceInternalError
	}
	if req == nil {
		return nil, status.StatusInternalParamsError
	}
	listOp := s.documentRepo.ListDocuments(&model.DocumentMeta{})
	if req.MetaArgs != nil {
		listOp = listOp.WithUserID(req.MetaArgs.UserID).
			WithParentID(req.MetaArgs.ParentID).
			WithKnowledgeID(req.MetaArgs.KnowledgeID)
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

func (s *DocumentService) listDocuments(req *documentsListReq) ([]*dto.DocumentResponse, int64, error) {
	listOp, err := s.buildListDocumentsOperation(req)
	if err != nil {
		return nil, 0, err
	}
	models, total, err := listOp.ExecWithTotal()

	if err != nil {
		return nil, 0, err
	}
	responses := make([]*dto.DocumentResponse, 0, len(models))
	for _, model := range models {
		responses = append(responses, s.ConvertToDocumentResponse(&model))
	}
	return responses, total, nil
}

func (s *DocumentService) getChildDocuments(parentID int64, options *dto.RecursiveOptions) ([]*dto.DocumentResponse, error) {
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

func (s *DocumentService) listDocumentsWithChildren(req *documentsListReq) ([]*dto.DocumentResponse, int64, error) {
	docs, total, err := s.listDocuments(req)
	if err != nil {
		return nil, 0, err
	}
	if req.Options != nil && req.Options.IncludeChildren {
		listOp, err := s.buildListDocumentsOperation(req)
		if err != nil {
			return nil, 0, err
		}
		models, _, err := listOp.ExecWithTotal()
		if err != nil {
			return nil, 0, err
		}

		for i := range docs {
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

func (s *DocumentService) ListDocumentsWithChildrenByExternal(req *dto.ListDocumentsByExternalReq) ([]*dto.DocumentResponse, int64, error) {
	var userID *int64
	if req.ExternalArgs.UserExternalID != nil && *req.ExternalArgs.UserExternalID != "" {
		id, err := repository.UserRepo.GetIDByExternalID(*req.ExternalArgs.UserExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		userID = &id
	}

	var parentID int64
	if req.ExternalArgs.RootDocumentExternalID != nil && *req.ExternalArgs.RootDocumentExternalID != "" {
		doc, err := s.GetDocumentByExternalID(*req.ExternalArgs.RootDocumentExternalID)
		if err != nil {
			return nil, 0, err
		}
		parentID = doc.ID
	}

	var knowledgeID *int64
	if req.ExternalArgs.KnowledgeExternalID != nil && *req.ExternalArgs.KnowledgeExternalID != "" {
		id, err := repository.KnowledgeBaseRepo.GetIDByExternalID(*req.ExternalArgs.KnowledgeExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		knowledgeID = &id
	}

	return s.listDocumentsWithChildren(
		&documentsListReq{
			MetaArgs: &documentsListMetaArgs{
				UserID:      userID,
				ParentID:    &parentID, // default to root
				KnowledgeID: knowledgeID,
			},
			FilterArgs: req.FilterArgs,
			SortArgs:   req.SortArgs,
			Pagination: req.Pagination,
			Options:    req.Options,
		},
	)
}

func (s *DocumentService) Create(req *dto.CreateDocumentReq) (*dto.CreateDocumentResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// TODO: redis support
	docExternalID := utils.GenerateExternalID(utils.ExternalIDContextDocument)
	err := utils.WithTransaction(func(tx *gorm.DB) error {
		docModel := &model.DocumentMeta{
			ExternalID: docExternalID,
			Title:      req.Title,
			Type:       req.Type,
			Summary:    req.Summary,
			Content:    req.Content,
		}
		// query user_id
		userID, err := s.userRepo.GetIDByExternalID(*req.CreatorExternalID).WithTx(tx).Exec()
		if err != nil {
			return status.StatusUserNotFound
		}

		docModel.UserID = userID
		// query knowledge_id
		var ownerID *int64
		if !req.CreateAsOwnDoc {
			kbID, err := s.kbRepo.GetIDByExternalID(*req.KnowledgeExternalID).WithTx(tx).Exec()
			if err != nil {
				return status.StatusKnowledgeBaseNotFound
			}
			docModel.KnowledgeID = &kbID
		} else {
			ownerID = &userID
		}
		// query parent_id
		if req.ParentExternalID != nil {
			parentDocID, err := s.documentRepo.GetIDByExternalID(*req.ParentExternalID).WithTx(tx).Exec()
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

		// create doc-kb relation
		err = s.docKBRelationRepo.CreateDocKBRelation(docModel.ID).
			WithKnowledgeID(docModel.KnowledgeID).
			WithOwnerID(ownerID).
			Exec()
		if err != nil {
			return status.GenErrWithCustomMsg(status.StatusWriteDBError, "create doc-kb relation failed")
		}

		// create doc version
		// err = s.
		// TODO:

		// create content

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
