package main

import (
	"github.com/sirupsen/logrus"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	logrus.Println("Starting migration...")

	if err := runMigration(); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
	logrus.Println("Database migration completed successfully")

}
