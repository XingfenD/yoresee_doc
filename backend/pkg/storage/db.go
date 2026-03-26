package storage

import (
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitPostgres(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	gormLogLevel, err := resolveGormLogLevel()
	if err != nil {
		return err
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	return nil
}

func resolveGormLogLevel() (logger.LogLevel, error) {
	levelText := ""
	if config.GlobalConfig != nil {
		levelText = strings.TrimSpace(config.GlobalConfig.Backend.Log.GormLogLevel)
		if levelText == "" {
			levelText = strings.TrimSpace(config.GlobalConfig.Backend.Log.Level)
		}
	}
	if levelText == "" {
		levelText = "info"
	}

	switch strings.ToLower(levelText) {
	case "silent":
		return logger.Silent, nil
	case "error":
		return logger.Error, nil
	case "warn", "warning":
		return logger.Warn, nil
	case "info", "debug", "trace":
		return logger.Info, nil
	case "fatal", "panic":
		return logger.Error, nil
	default:
		return logger.Info, fmt.Errorf("invalid gorm log level: %s", levelText)
	}
}

func ClosePostgres() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// func InitMySQL(cfg *config.DatabaseConfig) error {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

// 	var err error
// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	sqlDB, err := DB.DB()
// 	if err != nil {
// 		return err
// 	}

// 	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
// 	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

// 	return nil
// }
