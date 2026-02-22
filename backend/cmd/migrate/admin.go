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

	if err := tx.Where(
		"resource_type = ? AND resource_id = ? AND subject_type = ? AND subject_id = ?",
		model.ResourceTypeNamespace, "default", model.SubjectTypeUser, adminUser.ExternalID,
	).Delete(&model.PermissionRule{}).Error; err != nil {
		return err
	}

	if err := tx.Exec(
		`INSERT INTO permission_rules (
			resource_type, resource_id, resource_path,
			subject_type, subject_id, permissions, scope_type,
			is_deny, priority, created_by, created_at
		) VALUES (
			?, ?, ?, ?, ?,
			ARRAY['read', 'edit', 'manage', 'admin', 'create', 'transfer', 'audit'],
			?, ?, ?, ?, NOW()
		)`,
		model.ResourceTypeNamespace, "default", "",
		model.SubjectTypeUser, adminUser.ExternalID,
		model.ScopeTypeRecursive, false, 1, "",
	).Error; err != nil {
		return err
	}
	logrus.Println("Admin permission granted successfully in transaction.")
	return nil
}
