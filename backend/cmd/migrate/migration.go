package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func runMigration() error {
	if err := storage.DB.Exec("CREATE EXTENSION IF NOT EXISTS ltree").Error; err != nil {
		return err
	}
	logrus.Println("ltree extension created successfully")

	err := storage.DB.AutoMigrate(
		&model.Attachment{},
		&model.Document{},
		&model.DocumentYjsSnapshot{},
		&model.DocumentVersion{},
		&model.Invitation{},
		&model.InvitationRecord{},
		&model.KnowledgeBase{},
		&model.RecentKnowledgeBase{},
		&model.RecentTemplate{},
		&model.MembershipRelation{},
		&model.OrgNodeMeta{},
		&model.Template{},
		&model.UserGroupMeta{},
		&model.User{},
	)
	if err != nil {
		return err
	}

	if err := migrateRecentKnowledgeBaseIndex(storage.DB); err != nil {
		return err
	}
	if err := migrateRecentTemplateIndex(storage.DB); err != nil {
		return err
	}

	return nil
}

func migrateRecentKnowledgeBaseIndex(db *gorm.DB) error {
	return db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_recent_kb_user_kb ON recent_knowledge_bases (user_id, knowledge_base_id)").Error
}

func migrateRecentTemplateIndex(db *gorm.DB) error {
	return db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_recent_tpl_user_tpl ON recent_templates (user_id, template_id)").Error
}
