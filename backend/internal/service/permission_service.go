package service

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo *repository.PermissionRepository
	namespaceRepo  *repository.NamespaceRepository
	resourceRepo   *repository.ResourceRepository
	subjectRepo    *repository.SubjectRepository
}

// NewPermissionService 创建权限服务实例
func NewPermissionService() *PermissionService {
	return &PermissionService{
		permissionRepo: repository.PermissionRepo,
		namespaceRepo:  repository.NamespaceRepo,
		resourceRepo:   repository.ResourceRepo,
		subjectRepo:    repository.SubjectRepo,
	}
}

var PermissionSvc = NewPermissionService()

// GrantPermission 授予权限
func (s *PermissionService) GrantPermission(grant *dto.PermissionGrant) error {
	rule := s.convertToPermissionRule(grant)
	return s.permissionRepo.CreateRule(rule).Exec()
}

// convertToPermissionRule 将PermissionGrant转换为PermissionRule
func (s *PermissionService) convertToPermissionRule(grant *dto.PermissionGrant) *model.PermissionRule {
	return &model.PermissionRule{
		Namespace:    grant.Namespace,
		ResourceType: model.ResourceType(grant.Resource.Type),
		ResourceID:   grant.Resource.ID,
		ResourcePath: grant.Resource.Path,
		SubjectType:  model.SubjectType(grant.Subject.Type),
		SubjectID:    grant.Subject.ID,
		Permissions:  grant.Permissions,
		ScopeType:    model.ScopeType(grant.Scope.Type),
		IsDeny:       grant.IsDeny,
		Priority:     grant.Priority,
		ValidFrom:    grant.ValidFrom,
		ValidUntil:   grant.ValidUntil,
	}
}

// RevokePermission 撤销权限
func (s *PermissionService) RevokePermission(id int) error {
	return s.permissionRepo.DeleteRule(id).Exec()
}

// GetPermissionRule 获取权限规则
func (s *PermissionService) GetPermissionRule(id int) (*model.PermissionRule, error) {
	return s.permissionRepo.GetRuleByID(id).Exec()
}

// CreateNamespace 创建命名域
func (s *PermissionService) CreateNamespace(create *dto.NamespaceCreate) error {
	namespace := s.convertToNamespace(create)
	return s.namespaceRepo.Create(namespace).Exec()
}

// convertToNamespace 将NamespaceCreate转换为Namespace
func (s *PermissionService) convertToNamespace(create *dto.NamespaceCreate) *model.Namespace {
	return &model.Namespace{
		ID:      create.ID,
		Name:    create.Name,
		OwnerID: create.OwnerID,
	}
}

// GetNamespace 获取命名域
func (s *PermissionService) GetNamespace(id string) (*dto.NamespaceResponse, error) {
	namespace, err := s.namespaceRepo.GetByID(id).Exec()
	if err != nil {
		return nil, err
	}
	return dto.NewNamespaceResponseFromModel(namespace), nil
}

// CreateResource 创建资源
func (s *PermissionService) CreateResource(resource *model.Resource) error {
	return s.resourceRepo.Create(resource).Exec()
}

// GetResource 获取资源
func (s *PermissionService) GetResource(id string) (*model.Resource, error) {
	return s.resourceRepo.GetByID(id).Exec()
}

// CreateSubject 创建主体
func (s *PermissionService) CreateSubject(subject *model.Subject) error {
	return s.subjectRepo.Create(subject).Exec()
}

// GetSubject 获取主体
func (s *PermissionService) GetSubject(id string) (*model.Subject, error) {
	return s.subjectRepo.GetByID(id).Exec()
}

// CheckPermission 检查用户是否对资源有权限
func (s *PermissionService) CheckPermission(userID int64, namespace string, check *dto.PermissionCheck) (bool, error) {
	// 1. 获取用户的所有主体身份
	subjects, err := s.getUserSubjects(userID, namespace)
	if err != nil {
		return false, err
	}

	// 2. 获取资源的所有相关权限规则
	rules, err := s.getRelevantRules(model.ResourceType(check.Resource.Type), check.Resource.ID, namespace, subjects)
	if err != nil {
		return false, err
	}

	// 3. 解析最终权限
	effectivePermissions := s.resolveEffectivePermissions(rules)

	// 4. 检查是否包含所需权限
	return s.hasPermission(effectivePermissions, check.Permission), nil
}

// getUserSubjects 获取用户的所有主体身份
func (s *PermissionService) getUserSubjects(userID int64, namespace string) ([]model.Subject, error) {
	// TODO: 实现完整的用户主体身份获取逻辑
	// 暂时只返回用户本身作为主体
	userSubject := model.Subject{
		ID:        fmt.Sprintf("%d", userID),
		Type:      model.SubjectTypeUser,
		Namespace: namespace,
	}
	return []model.Subject{userSubject}, nil
}

// getRelevantRules 获取资源的所有相关权限规则
func (s *PermissionService) getRelevantRules(resourceType model.ResourceType, resourceID string, namespace string, subjects []model.Subject) ([]model.PermissionRule, error) {
	// TODO: 实现完整的权限规则获取逻辑
	// 暂时返回空数组，表示没有权限规则
	return []model.PermissionRule{}, nil
}

// resolveEffectivePermissions 解析最终权限
func (s *PermissionService) resolveEffectivePermissions(rules []model.PermissionRule) []string {
	// TODO: 实现完整的权限解析逻辑
	// 暂时返回空数组，表示没有权限
	return []string{}
}

// hasPermission 检查是否包含所需权限
func (s *PermissionService) hasPermission(permissions []string, requiredPermission string) bool {
	for _, perm := range permissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}

// BatchCheckPermissions 批量检查权限
func (s *PermissionService) BatchCheckPermissions(userID int64, namespace string, batchCheck *dto.PermissionBatchCheck) (dto.PermissionBatchCheckResponse, error) {
	// 1. 获取用户的所有主体身份
	subjects, err := s.getUserSubjects(userID, namespace)
	if err != nil {
		return nil, err
	}

	// 2. 获取资源的所有相关权限规则
	rules, err := s.getRelevantRules(model.ResourceType(batchCheck.Resource.Type), batchCheck.Resource.ID, namespace, subjects)
	if err != nil {
		return nil, err
	}

	// 3. 解析最终权限
	effectivePermissions := s.resolveEffectivePermissions(rules)

	// 4. 批量检查权限
	result := make(dto.PermissionBatchCheckResponse)
	for _, perm := range batchCheck.Permissions {
		result[perm] = s.hasPermission(effectivePermissions, perm)
	}

	return result, nil
}
