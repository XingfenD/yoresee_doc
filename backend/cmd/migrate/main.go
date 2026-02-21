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
	err := storage.DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.DocumentVersion{},
		&model.KnowledgeBase{},
		&model.RecentKnowledgeBase{},
		&model.Invitation{},
		&model.SystemConfig{},
		&model.DocumentCollaborator{},
	)

	if err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
}
