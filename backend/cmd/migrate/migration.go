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
		&model.Attachment{},
		&model.Content{},
		&model.DocumentMeta{},
		&model.DocumentVersion{},
		&model.Invitation{},
		&model.KnowledgeBase{},
		&model.RecentKnowledgeBase{},
		&model.DocKnowledgeRelation{},
		&model.MembershipRelation{},
		&model.OrgNodeMeta{},
		&model.PermissionRule{},
		&model.Resource{},
		&model.SystemConfig{},
		&model.TemplateMeta{},
		&model.UserGroupMeta{},
		&model.User{},
	)

	return err
}
