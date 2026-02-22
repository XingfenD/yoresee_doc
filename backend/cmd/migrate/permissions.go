package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func initializePermissions() error {
	logrus.Println("Starting permission initialization...")

	// 1. 创建默认命名域
	defaultNamespace := model.Namespace{
		ID:      "default",
		Name:    "Default Namespace",
		OwnerID: "admin",
	}
	// 使用Upsert操作，避免重复键错误
	if err := storage.DB.Where("id = ?", defaultNamespace.ID).FirstOrCreate(&defaultNamespace).Error; err != nil {
		return err
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
		return err
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
		return err
	}
	logrus.Println("Namespace resource created successfully.")

	logrus.Println("Permission initialization completed successfully.")
	return nil
}
