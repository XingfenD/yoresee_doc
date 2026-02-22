package service

import (
	"context"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

// DocumentService 文档服务
type DocumentService struct {
	documentRepo  *repository.DocumentRepository
	permissionSvc *PermissionService
}

// NewDocumentService 创建文档服务实例
func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo:  repository.DocumentRepo,
		permissionSvc: PermissionSvc,
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

// GetDocumentContent 获取文档内容
func (s *DocumentService) GetDocumentContent(documentID int64) (string, error) {
	return s.documentRepo.GetContent(documentID).Exec()
}

// CheckDocumentPermission 检查文档权限
func (s *DocumentService) CheckDocumentPermission(userID int64, documentID int64, namespace string, requiredPermission string) (bool, error) {
	cacheKey := fmt.Sprintf("permission:user:%d:doc:%d:%s", userID, documentID, requiredPermission)
	ctx := context.Background()

	cachedResult, err := storage.GetCache(ctx, cacheKey)
	if err == nil {
		return cachedResult == "true", nil
	}

	permissionCheck := &dto.PermissionCheck{
		Resource: dto.Resource{
			Type: string(model.ResourceTypeDocument),
			ID:   fmt.Sprintf("%d", documentID),
		},
		Permission: requiredPermission,
	}
	allowed, err := s.permissionSvc.CheckPermission(userID, namespace, permissionCheck)
	if err != nil {
		return false, err
	}

	err = storage.SetCache(ctx, cacheKey, fmt.Sprintf("%v", allowed), time.Hour)
	if err != nil {
		fmt.Printf("Set permission cache failed: %v\n", err)
	}

	return allowed, nil
}

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
