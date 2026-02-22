package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	// Create Admin User
	password := "admin123456"
	hashedPwd, _ := utils.HashPassword(password)
	externalID := utils.GenerateExternalID("admin")

	adminUser := model.User{
		ExternalID:   externalID,
		Username:     "admin",
		PasswordHash: hashedPwd,
		Email:        "admin@example.com",
		Nickname:     "Admin",
		Status:       1,
	}

	var count int64
	storage.DB.Model(&model.User{}).Where("email = ?", adminUser.Email).Count(&count)
	if count == 0 {
		if err := storage.DB.Create(&adminUser).Error; err != nil {
			logrus.Fatalf("Create admin user failed: %v", err)
		}
		logrus.Println("Admin user created successfully.")
	} else {
		logrus.Println("Admin user already exists.")
	}

	// 4. 授予系统管理员所有权限
	// 使用原始SQL语句，确保正确处理PostgreSQL数组类型
	// 先删除可能存在的记录
	if err := storage.DB.Exec(
		"DELETE FROM permission_rules WHERE namespace = ? AND resource_type = ? AND resource_id = ? AND subject_type = ? AND subject_id = ?",
		"default", model.ResourceTypeNamespace, "default", model.SubjectTypeUser, adminUser.ExternalID,
	).Error; err != nil {
		logrus.Fatalf("Delete existing permission rule failed: %v", err)
	}

	// 使用正确的PostgreSQL数组语法插入权限
	if err := storage.DB.Exec(
		`INSERT INTO permission_rules (
			namespace, resource_type, resource_id, resource_path, 
			subject_type, subject_id, permissions, scope_type, 
			is_deny, priority, created_by, created_at
		) VALUES (
			?, ?, ?, ?, ?, ?, 
			ARRAY['read', 'edit', 'manage', 'admin', 'create', 'transfer', 'audit'], 
			?, ?, ?, ?, NOW()
		)`,
		"default", model.ResourceTypeNamespace, "default", "",
		model.SubjectTypeUser, adminUser.ExternalID,
		model.ScopeTypeRecursive, false, 1, "",
	).Error; err != nil {
		logrus.Fatalf("Grant admin permission failed: %v", err)
	}
	logrus.Println("Admin permission granted successfully.")
}
