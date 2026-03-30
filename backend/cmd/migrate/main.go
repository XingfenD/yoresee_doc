package main

import (
	"github.com/sirupsen/logrus"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
)

func main() {
	initializer := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres()
	if err := initializer.Err(); err != nil {
		logrus.Fatalf("Init migrate failed: %v", err)
	}
	defer initializer.Shutdown()

	logrus.Println("Starting migration...")

	if err := runMigration(); err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}
	logrus.Println("Database migration completed successfully")

}
