package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	// 初始化数据库连接
	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	// 检查数据库是否已初始化
	if isDatabaseInitialized() {
		logrus.Println("Database already initialized, skipping initialization steps")
	} else {
		// 在事务中执行初始化操作
		if err := initializeDatabaseInTransaction(); err != nil {
			logrus.Fatalf("Database initialization failed: %v", err)
		}
		logrus.Println("Database initialized successfully")
	}

	logrus.Println("All initialization tasks completed successfully!")
}

func isDatabaseInitialized() bool {
	// 直接查询数据库，避免 Redis 依赖
	var config model.SystemConfig
	err := storage.DB.Where("key = ?", utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Database,
		constant.ConfigKey_Third_Initialized,
	)).First(&config).Error
	if err != nil {
		return false
	}
	return config.Value == constant.Database_Initialized_True
}

func initializeDatabaseInTransaction() error {
	return utils.WithTransaction(func(tx *gorm.DB) error {
		// 执行各个初始化函数
		if err := initializeConfigInTx(tx); err != nil {
			return err
		}

		if err := initializePermissionsInTx(tx); err != nil {
			return err
		}

		if err := createAdminUserInTx(tx); err != nil {
			return err
		}

		if err := initializeDocumentsInTx(tx); err != nil {
			return err
		}

		// 标记数据库已初始化
		initializedConfig := &model.SystemConfig{
			Key: utils.GenConfigKey(
				constant.ConfigKey_First_System,
				constant.ConfigKey_Second_Database,
				constant.ConfigKey_Third_Initialized,
			),
			Value: constant.Database_Initialized_True,
		}
		if err := tx.FirstOrCreate(initializedConfig, model.SystemConfig{
			Key: initializedConfig.Key,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}
