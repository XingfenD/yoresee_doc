package main

import (
	"log"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		log.Fatalf("Init Postgres failed: %v", err)
	}

	log.Println("Starting migration...")
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
		log.Fatalf("Migration failed: %v", err)
	}
}
