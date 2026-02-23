package dto

import (
	"time"
)

// PermissionGrant 权限授予请求
type PermissionGrant struct {
	Resource    Resource   `json:"resource" binding:"required"`
	Subject     Subject    `json:"subject" binding:"required"`
	Permissions []string   `json:"permissions" binding:"required"`
	Scope       Scope      `json:"scope" binding:"required"`
	Priority    int        `json:"priority"`
	IsDeny      bool       `json:"is_deny"`
	ValidFrom   *time.Time `json:"valid_from"`
	ValidUntil  *time.Time `json:"valid_until"`
}

// Resource 资源DTO
type Resource struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
	Path string `json:"path"`
}
type Subject struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

type Scope struct {
	Type string `json:"type" binding:"required"`
}

type PermissionCheck struct {
	Resource   Resource `json:"resource" binding:"required"`
	Permission string   `json:"permission" binding:"required"`
}

type PermissionBatchCheck struct {
	Resource    Resource `json:"resource" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

type PermissionBatchCheckResponse map[string]bool

type PermissionEffectiveResponse []string
