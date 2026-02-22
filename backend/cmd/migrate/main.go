package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
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

	// Run database migration
	if err := runMigration(); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
	logrus.Println("Database migration completed successfully")

	// Initialize system config
	if err := initializeConfig(); err != nil {
		logrus.Fatalf("Initialize config failed: %v", err)
	}

	// Initialize permissions
	if err := initializePermissions(); err != nil {
		logrus.Fatalf("Initialize permissions failed: %v", err)
	}

	// Create admin user
	if err := createAdminUser(); err != nil {
		logrus.Fatalf("Create admin user failed: %v", err)
	}

	// Initialize documents
	if err := initializeDocuments(); err != nil {
		logrus.Fatalf("Initialize documents failed: %v", err)
	}

	logrus.Println("All migration tasks completed successfully!")
}
