package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

// PermissionRepository 权限仓库
type PermissionRepository struct{}

var PermissionRepo = &PermissionRepository{}

// PermissionRuleCreateOperation 权限规则创建操作
type PermissionRuleCreateOperation struct {
	repo *PermissionRepository
	rule *model.PermissionRule
	tx   *gorm.DB
}

func (r *PermissionRepository) CreateRule(rule *model.PermissionRule) *PermissionRuleCreateOperation {
	return &PermissionRuleCreateOperation{
		repo: r,
		rule: rule,
	}
}

func (op *PermissionRuleCreateOperation) WithTx(tx *gorm.DB) *PermissionRuleCreateOperation {
	op.tx = tx
	return op
}

func (op *PermissionRuleCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.rule).Error
	}
	return storage.DB.Create(op.rule).Error
}

// PermissionRuleGetByIDOperation 权限规则根据ID获取操作
type PermissionRuleGetByIDOperation struct {
	repo *PermissionRepository
	id   int
	tx   *gorm.DB
}

func (r *PermissionRepository) GetRuleByID(id int) *PermissionRuleGetByIDOperation {
	return &PermissionRuleGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *PermissionRuleGetByIDOperation) WithTx(tx *gorm.DB) *PermissionRuleGetByIDOperation {
	op.tx = tx
	return op
}

func (op *PermissionRuleGetByIDOperation) Exec() (*model.PermissionRule, error) {
	var rule model.PermissionRule
	var err error

	if op.tx != nil {
		err = op.tx.First(&rule, op.id).Error
	} else {
		err = storage.DB.First(&rule, op.id).Error
	}

	return &rule, err
}

// PermissionRuleDeleteOperation 权限规则删除操作
type PermissionRuleDeleteOperation struct {
	repo *PermissionRepository
	id   int
	tx   *gorm.DB
}

func (r *PermissionRepository) DeleteRule(id int) *PermissionRuleDeleteOperation {
	return &PermissionRuleDeleteOperation{
		repo: r,
		id:   id,
	}
}

func (op *PermissionRuleDeleteOperation) WithTx(tx *gorm.DB) *PermissionRuleDeleteOperation {
	op.tx = tx
	return op
}

func (op *PermissionRuleDeleteOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Delete(&model.PermissionRule{}, op.id).Error
	}
	return storage.DB.Delete(&model.PermissionRule{}, op.id).Error
}

// PermissionRuleGetByResourceOperation 根据资源获取权限规则操作
type PermissionRuleGetByResourceOperation struct {
	repo         *PermissionRepository
	namespace    string
	resourceType model.ResourceType
	resourceID   string
	tx           *gorm.DB
}

func (r *PermissionRepository) GetRulesByResource(namespace string, resourceType model.ResourceType, resourceID string) *PermissionRuleGetByResourceOperation {
	return &PermissionRuleGetByResourceOperation{
		repo:         r,
		namespace:    namespace,
		resourceType: resourceType,
		resourceID:   resourceID,
	}
}

func (op *PermissionRuleGetByResourceOperation) WithTx(tx *gorm.DB) *PermissionRuleGetByResourceOperation {
	op.tx = tx
	return op
}

func (op *PermissionRuleGetByResourceOperation) Exec() ([]model.PermissionRule, error) {
	var rules []model.PermissionRule
	var err error

	if op.tx != nil {
		err = op.tx.Where("namespace = ? AND resource_type = ? AND resource_id = ?", op.namespace, op.resourceType, op.resourceID).Find(&rules).Error
	} else {
		err = storage.DB.Where("namespace = ? AND resource_type = ? AND resource_id = ?", op.namespace, op.resourceType, op.resourceID).Find(&rules).Error
	}

	return rules, err
}

// PermissionRuleGetBySubjectOperation 根据主体获取权限规则操作
type PermissionRuleGetBySubjectOperation struct {
	repo        *PermissionRepository
	namespace   string
	subjectType model.SubjectType
	subjectID   string
	tx          *gorm.DB
}

func (r *PermissionRepository) GetRulesBySubject(namespace string, subjectType model.SubjectType, subjectID string) *PermissionRuleGetBySubjectOperation {
	return &PermissionRuleGetBySubjectOperation{
		repo:        r,
		namespace:   namespace,
		subjectType: subjectType,
		subjectID:   subjectID,
	}
}

func (op *PermissionRuleGetBySubjectOperation) WithTx(tx *gorm.DB) *PermissionRuleGetBySubjectOperation {
	op.tx = tx
	return op
}

func (op *PermissionRuleGetBySubjectOperation) Exec() ([]model.PermissionRule, error) {
	var rules []model.PermissionRule
	var err error

	if op.tx != nil {
		err = op.tx.Where("namespace = ? AND subject_type = ? AND subject_id = ?", op.namespace, op.subjectType, op.subjectID).Find(&rules).Error
	} else {
		err = storage.DB.Where("namespace = ? AND subject_type = ? AND subject_id = ?", op.namespace, op.subjectType, op.subjectID).Find(&rules).Error
	}

	return rules, err
}

// NamespaceRepository 命名域仓库
type NamespaceRepository struct{}

var NamespaceRepo = &NamespaceRepository{}

// NamespaceCreateOperation 命名域创建操作
type NamespaceCreateOperation struct {
	repo      *NamespaceRepository
	namespace *model.Namespace
	tx        *gorm.DB
}

func (r *NamespaceRepository) Create(namespace *model.Namespace) *NamespaceCreateOperation {
	return &NamespaceCreateOperation{
		repo:      r,
		namespace: namespace,
	}
}

func (op *NamespaceCreateOperation) WithTx(tx *gorm.DB) *NamespaceCreateOperation {
	op.tx = tx
	return op
}

func (op *NamespaceCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.namespace).Error
	}
	return storage.DB.Create(op.namespace).Error
}

// NamespaceGetByIDOperation 命名域根据ID获取操作
type NamespaceGetByIDOperation struct {
	repo *NamespaceRepository
	id   string
	tx   *gorm.DB
}

func (r *NamespaceRepository) GetByID(id string) *NamespaceGetByIDOperation {
	return &NamespaceGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *NamespaceGetByIDOperation) WithTx(tx *gorm.DB) *NamespaceGetByIDOperation {
	op.tx = tx
	return op
}

func (op *NamespaceGetByIDOperation) Exec() (*model.Namespace, error) {
	var namespace model.Namespace
	var err error

	if op.tx != nil {
		err = op.tx.First(&namespace, "id = ?", op.id).Error
	} else {
		err = storage.DB.First(&namespace, "id = ?", op.id).Error
	}

	return &namespace, err
}

// ResourceRepository 资源仓库
type ResourceRepository struct{}

var ResourceRepo = &ResourceRepository{}

// ResourceCreateOperation 资源创建操作
type ResourceCreateOperation struct {
	repo     *ResourceRepository
	resource *model.Resource
	tx       *gorm.DB
}

func (r *ResourceRepository) Create(resource *model.Resource) *ResourceCreateOperation {
	return &ResourceCreateOperation{
		repo:     r,
		resource: resource,
	}
}

func (op *ResourceCreateOperation) WithTx(tx *gorm.DB) *ResourceCreateOperation {
	op.tx = tx
	return op
}

func (op *ResourceCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.resource).Error
	}
	return storage.DB.Create(op.resource).Error
}

// ResourceGetByIDOperation 资源根据ID获取操作
type ResourceGetByIDOperation struct {
	repo *ResourceRepository
	id   string
	tx   *gorm.DB
}

func (r *ResourceRepository) GetByID(id string) *ResourceGetByIDOperation {
	return &ResourceGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *ResourceGetByIDOperation) WithTx(tx *gorm.DB) *ResourceGetByIDOperation {
	op.tx = tx
	return op
}

func (op *ResourceGetByIDOperation) Exec() (*model.Resource, error) {
	var resource model.Resource
	var err error

	if op.tx != nil {
		err = op.tx.First(&resource, "id = ?", op.id).Error
	} else {
		err = storage.DB.First(&resource, "id = ?", op.id).Error
	}

	return &resource, err
}

// SubjectRepository 主体仓库
type SubjectRepository struct{}

var SubjectRepo = &SubjectRepository{}

// SubjectCreateOperation 主体创建操作
type SubjectCreateOperation struct {
	repo    *SubjectRepository
	subject *model.Subject
	tx      *gorm.DB
}

func (r *SubjectRepository) Create(subject *model.Subject) *SubjectCreateOperation {
	return &SubjectCreateOperation{
		repo:    r,
		subject: subject,
	}
}

func (op *SubjectCreateOperation) WithTx(tx *gorm.DB) *SubjectCreateOperation {
	op.tx = tx
	return op
}

func (op *SubjectCreateOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.subject).Error
	}
	return storage.DB.Create(op.subject).Error
}

// SubjectGetByIDOperation 主体根据ID获取操作
type SubjectGetByIDOperation struct {
	repo *SubjectRepository
	id   string
	tx   *gorm.DB
}

func (r *SubjectRepository) GetByID(id string) *SubjectGetByIDOperation {
	return &SubjectGetByIDOperation{
		repo: r,
		id:   id,
	}
}

func (op *SubjectGetByIDOperation) WithTx(tx *gorm.DB) *SubjectGetByIDOperation {
	op.tx = tx
	return op
}

func (op *SubjectGetByIDOperation) Exec() (*model.Subject, error) {
	var subject model.Subject
	var err error

	if op.tx != nil {
		err = op.tx.First(&subject, "id = ?", op.id).Error
	} else {
		err = storage.DB.First(&subject, "id = ?", op.id).Error
	}

	return &subject, err
}
