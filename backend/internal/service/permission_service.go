package service

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

// PermissionService Permission service

type PermissionService struct {
	permissionRepo *repository.PermissionRepository
	resourceRepo   *repository.ResourceRepository
	subjectRepo    *repository.SubjectRepository
}

// NewPermissionService Create permission service instance
func NewPermissionService() *PermissionService {
	return &PermissionService{
		permissionRepo: repository.PermissionRepo,
		resourceRepo:   repository.ResourceRepo,
		subjectRepo:    repository.SubjectRepo,
	}
}

var PermissionSvc = NewPermissionService()

// GrantPermission Grant permission
func (s *PermissionService) GrantPermission(grant *dto.PermissionGrant) error {
	rule := s.convertToPermissionRule(grant)
	return s.permissionRepo.CreateRule(rule).Exec()
}

// convertToPermissionRule Convert PermissionGrant to PermissionRule
func (s *PermissionService) convertToPermissionRule(grant *dto.PermissionGrant) *model.PermissionRule {
	return &model.PermissionRule{
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

// RevokePermission Revoke permission
func (s *PermissionService) RevokePermission(id int) error {
	return s.permissionRepo.DeleteRule(id).Exec()
}

// GetPermissionRule Get permission rule
func (s *PermissionService) GetPermissionRule(id int) (*model.PermissionRule, error) {
	return s.permissionRepo.GetRuleByID(id).Exec()
}

// CreateResource Create resource
func (s *PermissionService) CreateResource(resource *model.Resource) error {
	return s.resourceRepo.Create(resource).Exec()
}

// GetResource Get resource
func (s *PermissionService) GetResource(id string) (*model.Resource, error) {
	return s.resourceRepo.GetByID(id).Exec()
}

// CreateSubject Create subject
func (s *PermissionService) CreateSubject(subject *model.Subject) error {
	return s.subjectRepo.Create(subject).Exec()
}

// GetSubject Get subject
func (s *PermissionService) GetSubject(id string) (*model.Subject, error) {
	return s.subjectRepo.GetByID(id).Exec()
}

// CheckPermission Check if user has permission on resource
func (s *PermissionService) CheckPermission(userID int64, check *dto.PermissionCheck) (bool, error) {
	subjects, err := s.getUserSubjects(userID)
	if err != nil {
		return false, err
	}

	rules, err := s.getRelevantRules(model.ResourceType(check.Resource.Type), check.Resource.ID, subjects)
	if err != nil {
		return false, err
	}

	effectivePermissions := s.resolveEffectivePermissions(rules)

	return s.hasPermission(effectivePermissions, check.Permission), nil
}

// getUserSubjects Get all subject identities for a user
func (s *PermissionService) getUserSubjects(userID int64) ([]model.Subject, error) {
	// Get user subject
	userSubject := model.Subject{
		ID:   fmt.Sprintf("%d", userID),
		Type: model.SubjectTypeUser,
	}

	// TODO: Add logic to get user groups and organization nodes
	// For now, only return the user itself as subject
	return []model.Subject{userSubject}, nil
}

// getRelevantRules Get all relevant permission rules for resources
func (s *PermissionService) getRelevantRules(resourceType model.ResourceType, resourceID string, subjects []model.Subject) ([]model.PermissionRule, error) {
	var allRules []model.PermissionRule

	// Get rules for the specific resource
	resourceRules, err := s.permissionRepo.GetRulesByResource(resourceType, resourceID).Exec()
	if err != nil {
		return nil, status.StatusReadDBError
	}
	allRules = append(allRules, resourceRules...)

	// Get rules for each subject
	for _, subject := range subjects {
		subjectRules, err := s.permissionRepo.GetRulesBySubject(subject.Type, subject.ID).Exec()
		if err != nil {
			return nil, status.StatusReadDBError
		}
		allRules = append(allRules, subjectRules...)
	}

	return allRules, nil
}

// resolveEffectivePermissions Resolve effective permissions from rules
func (s *PermissionService) resolveEffectivePermissions(rules []model.PermissionRule) []string {
	// Map to store permissions with their priority
	permissionMap := make(map[string]int)

	// Process each rule
	for _, rule := range rules {
		for _, perm := range rule.Permissions {
			// Check if this permission already exists
			if existingPriority, exists := permissionMap[perm]; !exists || rule.Priority < existingPriority {
				// Add or update permission with higher priority (lower number)
				permissionMap[perm] = rule.Priority
			}
		}
	}

	// Convert map to slice
	permissions := make([]string, 0, len(permissionMap))
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions
}

// hasPermission Check if permissions contain the required permission
func (s *PermissionService) hasPermission(permissions []string, requiredPermission string) bool {
	for _, perm := range permissions {
		if perm == requiredPermission {
			return true
		}
	}
	return false
}

// BatchCheckPermissions Batch check permissions
func (s *PermissionService) BatchCheckPermissions(userID int64, batchCheck *dto.PermissionBatchCheck) (dto.PermissionBatchCheckResponse, error) {
	// 1. Get user subjects
	subjects, err := s.getUserSubjects(userID)
	if err != nil {
		return nil, err
	}

	// 2. Get relevant permission rules
	rules, err := s.getRelevantRules(model.ResourceType(batchCheck.Resource.Type), batchCheck.Resource.ID, subjects)
	if err != nil {
		return nil, err
	}

	// 3. Resolve effective permissions
	effectivePermissions := s.resolveEffectivePermissions(rules)

	// 4. Batch check permissions
	result := make(dto.PermissionBatchCheckResponse)
	for _, perm := range batchCheck.Permissions {
		result[perm] = s.hasPermission(effectivePermissions, perm)
	}

	return result, nil
}
