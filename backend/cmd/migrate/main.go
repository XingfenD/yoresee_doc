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

	logrus.Println("Starting migration...")

	// 创建ltree扩展
	if err := storage.DB.Exec("CREATE EXTENSION IF NOT EXISTS ltree").Error; err != nil {
		logrus.Fatalf("Create ltree extension failed: %v", err)
	}
	logrus.Println("ltree extension created successfully")

	// 执行自动迁移
	err := storage.DB.AutoMigrate(
		&model.User{},
		&model.DocumentVersion{},
		&model.KnowledgeBase{},
		&model.RecentKnowledgeBase{},
		&model.Invitation{},
		&model.SystemConfig{},
		&model.DocumentCollaborator{},
		&model.Namespace{},
		&model.Resource{},
		&model.Subject{},
		&model.PermissionRule{},
	)

	if err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
}
