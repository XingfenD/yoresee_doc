package main

import (
	"log"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		log.Fatalf("Init Postgres failed: %v", err)
	}

	adminRole := model.Role{
		ID:          1,
		Name:        "超级管理员",
		Code:        "admin",
		Description: "系统超级管理员",
		IsSystem:    true,
		Permissions: []string{"*"},
	}
	storage.DB.FirstOrCreate(&adminRole, model.Role{ID: 1})

	userRole := model.Role{
		ID:          2,
		Name:        "普通用户",
		Code:        "user",
		Description: "普通用户",
		IsSystem:    true,
		Permissions: []string{"read", "write"},
	}
	storage.DB.FirstOrCreate(&userRole, model.Role{ID: 2})

	// Create Admin User
	password := "admin123456"
	hashedPwd, _ := utils.HashPassword(password)

	adminUser := model.User{
		Username:     "admin",
		PasswordHash: hashedPwd,
		Email:        "admin@example.com",
		Nickname:     "Admin",
		RoleID:       1,
		Status:       1,
	}

	var count int64
	storage.DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		if err := storage.DB.Create(&adminUser).Error; err != nil {
			log.Fatalf("Create admin user failed: %v", err)
		}
		log.Println("Admin user created successfully.")
	} else {
		log.Println("Admin user already exists.")
	}
}
