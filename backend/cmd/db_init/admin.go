package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createAdminUserInTx(tx *gorm.DB) error {
	logrus.Println("Creating admin user in transaction...")

	password := "admin123456"
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
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
	tx.Model(&model.User{}).Where("email = ?", adminUser.Email).Count(&count)
	if count == 0 {
		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}
		logrus.Println("Admin user created successfully in transaction.")
	} else {
		logrus.Println("Admin user already exists in transaction.")
	}

	// 为管理员创建文档级别的权限规则，确保能访问所有文档
	permissionsString := "read,edit,manage,admin,create,transfer,audit"

	if err := tx.Exec(`
		INSERT INTO permission_rules (
			resource_type, resource_id, resource_path, 
			subject_type, subject_id, permissions, scope_type,
			is_deny, priority, created_by, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`,
		model.ResourceTypeDocument, "*", "", // 使用通配符*表示所有文档
		model.SubjectTypeUser, adminUser.ExternalID,
		permissionsString, model.ScopeTypeRecursive, false, 1, "",
	).Error; err != nil {
		return err
	}

	logrus.Println("Admin permission granted successfully in transaction.")
	return nil
}
