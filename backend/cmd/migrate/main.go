package main

import (
	"github.com/sirupsen/logrus"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func main() {
	if err := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		Err(); err != nil {
		logrus.Fatalf("Init migrate failed: %v", err)
	}
	defer storage.ClosePostgres()

	logrus.Println("Starting migration...")

	if err := runMigration(); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
	logrus.Println("Database migration completed successfully")

}
