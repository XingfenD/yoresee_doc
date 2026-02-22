package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func initializePermissions() error {
	logrus.Println("Starting permission initialization...")

	// 创建系统管理员主体
	adminSubject := model.Subject{
		ID:   "admin",
		Type: model.SubjectTypeUser,
	}
	// 使用Upsert操作，避免重复键错误
	if err := storage.DB.Where("id = ?", adminSubject.ID).FirstOrCreate(&adminSubject).Error; err != nil {
		return err
	}
	logrus.Println("Admin subject created successfully.")

	logrus.Println("Permission initialization completed successfully.")
	return nil
}
