package service

import (
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
)

type DocumentService struct {
	documentRepo *repository.DocumentRepository
	// permissionSvc *PermissionService
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo: repository.DocumentRepo,
		// permissionSvc: PermissionSvc,
	}
}

var DocumentSvc = NewDocumentService()

// ConvertToDocumentResponse 将DocumentMeta转换为DocumentResponse
func (s *DocumentService) ConvertToDocumentResponse(doc *model.DocumentMeta) *dto.DocumentResponse {
	return dto.NewDocumentResponseFromModel(doc)
}

// GetDocumentByExternalID 根据ExternalID获取文档
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

func (s *DocumentService) ListDocuments(
	userID *int64,
	parentID *int64,
	titleKeyword *string,
	docType *string,
	status *int,
	tags []string,
	createTimeRangeStart *string,
	createTimeRangeEnd *string,
	updateTimeRangeStart *string,
	updateTimeRangeEnd *string,
	sortField string,
	sortDesc bool,
	page int,
	pageSize int,
) ([]*dto.DocumentResponse, int64, error) {

	models, total, err := s.documentRepo.ListDocuments(&model.DocumentMeta{}).
		WithUserID(userID).
		WithParentID(parentID).
		WithTitleKeyword(titleKeyword).
		WithType(docType).
		WithStatus(status).
		WithTags(tags).
		WithCreateTimeRange(createTimeRangeStart, createTimeRangeEnd).
		WithUpdateTimeRange(updateTimeRangeStart, updateTimeRangeEnd).
		WithSort(sortField, sortDesc).
		WithPagination(page, pageSize).
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

// ListDocumentsOptions 用于文档查询的选项
type ListDocumentsOptions struct {
	IncludeChildren bool `json:"include_children"`
	Recursive       bool `json:"recursive"`
	Depth           int  `json:"depth"`
}

// GetChildDocuments 获取子文档
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

		// 如果启用递归并且深度允许
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

// ListDocumentsWithChildren 根据条件查询文档列表并递归获取子文档
func (s *DocumentService) ListDocumentsWithChildren(
	userID *int64,
	parentID *int64,
	titleKeyword *string,
	docType *string,
	status *int,
	tags []string,
	createTimeRangeStart *string,
	createTimeRangeEnd *string,
	updateTimeRangeStart *string,
	updateTimeRangeEnd *string,
	sortField string,
	sortDesc bool,
	page int,
	pageSize int,
	options *ListDocumentsOptions,
) ([]*dto.DocumentResponse, int64, error) {

	// 首先获取文档列表
	docs, total, err := s.ListDocuments(
		userID,
		parentID,
		titleKeyword,
		docType,
		status,
		tags,
		createTimeRangeStart,
		createTimeRangeEnd,
		updateTimeRangeStart,
		updateTimeRangeEnd,
		sortField,
		sortDesc,
		page,
		pageSize,
	)
	if err != nil {
		return nil, 0, err
	}

	// 如果需要包含子文档，我们需要原始模型的ID
	if options != nil && options.IncludeChildren {
		// 重新获取模型以获得ID信息，用于递归查询子文档
		models, _, err := s.documentRepo.ListDocuments(&model.DocumentMeta{}).
			WithUserID(userID).
			WithParentID(parentID).
			WithTitleKeyword(titleKeyword).
			WithType(docType).
			WithStatus(status).
			WithTags(tags).
			WithCreateTimeRange(createTimeRangeStart, createTimeRangeEnd).
			WithUpdateTimeRange(updateTimeRangeStart, updateTimeRangeEnd).
			WithSort(sortField, sortDesc).
			WithPagination(page, pageSize).
			ExecWithTotal()
		if err != nil {
			return nil, 0, err
		}

		// 现在我们可以使用模型的ID来获取子文档
		for i := range docs {
			childDocs, err := s.GetChildDocuments(models[i].ID, options)
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

// ListDocumentsWithChildrenByExternalID 根据外部ID查询文档列表并递归获取子文档
func (s *DocumentService) ListDocumentsWithChildrenByExternalID(
	userExternalID *string,
	rootDocumentExternalID *string,
	titleKeyword *string,
	docType *string,
	status *int,
	tags []string,
	createTimeRangeStart *string,
	createTimeRangeEnd *string,
	updateTimeRangeStart *string,
	updateTimeRangeEnd *string,
	sortField string,
	sortDesc bool,
	page int,
	pageSize int,
	options *ListDocumentsOptions,
) ([]*dto.DocumentResponse, int64, error) {

	var userID *int64
	if userExternalID != nil && *userExternalID != "" {
		id, err := repository.UserRepo.GetIDByExternalID(*userExternalID).Exec()
		if err != nil {
			return nil, 0, err
		}
		userID = &id
	}

	var parentID *int64
	if rootDocumentExternalID != nil && *rootDocumentExternalID != "" {
		doc, err := s.GetDocumentByExternalID(*rootDocumentExternalID)
		if err != nil {
			return nil, 0, err
		}
		parentID = &doc.ID
	}

	return s.ListDocumentsWithChildren(
		userID,
		parentID,
		titleKeyword,
		docType,
		status,
		tags,
		createTimeRangeStart,
		createTimeRangeEnd,
		updateTimeRangeStart,
		updateTimeRangeEnd,
		sortField,
		sortDesc,
		page,
		pageSize,
		options,
	)
}
