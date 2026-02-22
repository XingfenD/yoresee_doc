package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func runMigration() error {
	if err := storage.DB.Exec("CREATE EXTENSION IF NOT EXISTS ltree").Error; err != nil {
		return err
	}
	logrus.Println("ltree extension created successfully")

	err := storage.DB.AutoMigrate(
		&model.User{},
		&model.DocumentMeta{},
		&model.DocumentVersion{},
		&model.Content{},
		&model.KnowledgeBase{},
		&model.RecentKnowledgeBase{},
		&model.Invitation{},
		&model.SystemConfig{},
		&model.Namespace{},
		&model.Resource{},
		&model.Subject{},
		&model.PermissionRule{},
	)

	return err
}
