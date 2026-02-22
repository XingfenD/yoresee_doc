package service

import (
	"context"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

// DocumentService 文档服务
type DocumentService struct {
	documentRepo  *repository.DocumentRepository
	permissionSvc *PermissionService
	userSvc       *UserService
}

// NewDocumentService 创建文档服务实例
func NewDocumentService() *DocumentService {
	return &DocumentService{
		documentRepo:  repository.DocumentRepo,
		permissionSvc: PermissionSvc,
		userSvc:       UserSvc,
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
	// 1. 生成缓存键
	cacheKey := fmt.Sprintf("permission:user:%d:doc:%d:%s", userID, documentID, requiredPermission)
	ctx := context.Background()

	// 2. 尝试从缓存获取权限检查结果
	cachedResult, err := storage.GetCache(ctx, cacheKey)
	if err == nil {
		// 缓存命中
		return cachedResult == "true", nil
	}

	// 3. 缓存未命中，执行权限检查
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

	// 4. 将权限检查结果存入缓存，有效期1小时
	err = storage.SetCache(ctx, cacheKey, fmt.Sprintf("%v", allowed), time.Hour)
	if err != nil {
		// 缓存失败不影响权限检查结果
		fmt.Printf("Set permission cache failed: %v\n", err)
	}

	return allowed, nil
}

// GetDocumentWithContent 根据ExternalID获取文档及其内容
func (s *DocumentService) GetDocumentWithContent(docExternalID string, userExternalID string, namespace string) (*model.DocumentMeta, string, error) {
	// 1. 获取文档元数据
	document, err := s.GetDocumentByExternalID(docExternalID)
	if err != nil {
		return nil, "", err
	}

	userID, err := s.userSvc.GetIDByExternalID(userExternalID)
	if err != nil {
		return nil, "", status.StatusUserNotFound
	}

	// 2. 检查权限
	allowed, err := s.CheckDocumentPermission(userID, document.ID, namespace, string(model.PermissionRead))
	if err != nil {
		return nil, "", err
	}
	if !allowed {
		return nil, "", fmt.Errorf("permission denied")
	}

	// 3. 获取文档内容
	content, err := s.GetDocumentContent(document.ID)
	if err != nil {
		return document, "", err
	}

	return document, content, nil
}
