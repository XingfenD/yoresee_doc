package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
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

	logrus.Println("Starting permission initialization...")

	// 1. 创建默认命名域
	defaultNamespace := model.Namespace{
		ID:      "default",
		Name:    "Default Namespace",
		OwnerID: "admin",
	}
	// 使用Upsert操作，避免重复键错误
	if err := storage.DB.Where("id = ?", defaultNamespace.ID).FirstOrCreate(&defaultNamespace).Error; err != nil {
		logrus.Fatalf("Create default namespace failed: %v", err)
	}
	logrus.Println("Default namespace created successfully.")

	// 2. 创建系统管理员主体
	adminSubject := model.Subject{
		ID:        "admin",
		Namespace: "default",
		Type:      model.SubjectTypeUser,
	}
	// 使用Upsert操作，避免重复键错误
	if err := storage.DB.Where("id = ?", adminSubject.ID).FirstOrCreate(&adminSubject).Error; err != nil {
		logrus.Fatalf("Create admin subject failed: %v", err)
	}
	logrus.Println("Admin subject created successfully.")

	// 3. 创建默认命名域资源
	namespaceResource := model.Resource{
		ID:        "default",
		Namespace: "default",
		Type:      model.ResourceTypeNamespace,
	}
	// 使用Upsert操作，避免重复键错误
	if err := storage.DB.Where("id = ?", namespaceResource.ID).FirstOrCreate(&namespaceResource).Error; err != nil {
		logrus.Fatalf("Create namespace resource failed: %v", err)
	}
	logrus.Println("Namespace resource created successfully.")

	logrus.Println("Permission initialization completed successfully.")
}
