package storage

import (
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/pkg/errs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func OpenPostgres(cfg *PostgresOptions) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errs.ErrPostgresOptionsNil
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	gormLogLevel, err := resolveGormLogLevel(cfg.GormLogLevel)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}

	return db, nil
}

func resolveGormLogLevel(levelText string) (logger.LogLevel, error) {
	levelText = strings.TrimSpace(levelText)
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
		return logger.Info, errs.Detail(errs.ErrInvalidGormLogLevel, levelText)
	}
}

func ClosePostgres() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
