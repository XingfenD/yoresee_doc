package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type PermissionRepository struct{}

var PermissionRepo = &PermissionRepository{}

type CreateRoleOperation struct {
	repo *PermissionRepository
	role *model.Role
	tx   *gorm.DB
}

func (r *PermissionRepository) CreateRole(role *model.Role) *CreateRoleOperation {
	return &CreateRoleOperation{
		repo: r,
		role: role,
	}
}

func (op *CreateRoleOperation) WithTx(tx *gorm.DB) *CreateRoleOperation {
	op.tx = tx
	return op
}

func (op *CreateRoleOperation) Exec() error {
	if op.tx != nil {
		return op.tx.Create(op.role).Error
	}
	return storage.DB.Create(op.role).Error
}
