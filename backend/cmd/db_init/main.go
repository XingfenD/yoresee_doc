package main

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	if err := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitConsul().
		RequireConsulEnabled().
		Err(); err != nil {
		logrus.Fatalf("Init db_init failed: %v", err)
	}

	if isDatabaseInitialized() {
		logrus.Println("Database already initialized, skipping initialization steps")
	} else {
		if err := initializeDatabaseInTransaction(); err != nil {
			logrus.Fatalf("Database initialization failed: %v", err)
		}
		if err := initializeConfigInConsul(context.Background()); err != nil {
			logrus.Fatalf("Consul config initialization failed: %v", err)
		}
		if err := markDatabaseInitializedInConsul(context.Background()); err != nil {
			logrus.Fatalf("Mark database initialized failed: %v", err)
		}
		logrus.Println("Database initialized successfully")
	}

	logrus.Println("All initialization tasks completed successfully!")
}

func isDatabaseInitialized() bool {
	value, ok, err := storage.Consul.Get(context.Background(), utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Database,
		constant.ConfigKey_Third_Initialized,
	))
	if err != nil || !ok {
		return false
	}
	return value == constant.Database_Initialized_True
}

func initializeDatabaseInTransaction() error {
	return utils.WithTransaction(func(tx *gorm.DB) error {
		if err := initializePermissionsInTx(tx); err != nil {
			return err
		}

		if err := createAdminUserInTx(tx); err != nil {
			return err
		}

		if err := createTestUserIntx(tx); err != nil {
			return err
		}

		if err := initializeDocumentsInTx(tx); err != nil {
			return err
		}

		if err := initializeKnowledgeBasesInTx(tx); err != nil {
			return err
		}

		if err := createNormalUserInTx(tx); err != nil {
			return err
		}

		if err := initializeUserGroupsInTx(tx); err != nil {
			return err
		}

		return nil
	})
}
