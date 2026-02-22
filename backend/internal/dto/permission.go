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

// Subject 主体DTO
type Subject struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

// Scope 权限范围DTO
type Scope struct {
	Type string `json:"type" binding:"required"`
}

// PermissionCheck 权限检查请求
type PermissionCheck struct {
	Resource   Resource `json:"resource" binding:"required"`
	Permission string   `json:"permission" binding:"required"`
}

// PermissionBatchCheck 批量权限检查请求
type PermissionBatchCheck struct {
	Resource    Resource `json:"resource" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

// PermissionBatchCheckResponse 批量权限检查响应
type PermissionBatchCheckResponse map[string]bool

// PermissionEffectiveResponse 有效权限响应
type PermissionEffectiveResponse []string
