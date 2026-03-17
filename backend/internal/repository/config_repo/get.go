package config_repo

import (
	"context"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ConfigGetOperation struct {
	repo *ConfigRepository
	key  string
	tx   *gorm.DB
}

func (r *ConfigRepository) Get(key string) *ConfigGetOperation {
	// var config model.SystemConfig
	// return &config, storage.DB.Where("key = ?", key).First(&config).Error
	return &ConfigGetOperation{
		repo: r,
		key:  key,
	}
}

func (op *ConfigGetOperation) WithTx(tx *gorm.DB) *ConfigGetOperation {
	op.tx = tx
	return op
}

func (op *ConfigGetOperation) query() (*model.SystemConfig, error) {
	var config model.SystemConfig
	if err := storage.DB.Where("key = ?", op.key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (op *ConfigGetOperation) Exec(ctx context.Context) (*model.SystemConfig, error) {

	if op.tx == nil {
		op.tx = storage.DB
	}
	configCacheKey := cache.KeySystemConfig(op.key)
	config, err := cache_loader.NewCacheLoadOperation[model.SystemConfig](&op.repo.Loader).
		WithDBLoader(func() (*model.SystemConfig, error) {
			return op.query()
		}).
		WithTTL(cacheExpiration).
		WithDefaultKeyAndParser(configCacheKey, nil).
		Exec(ctx)

	if err != nil {
		logrus.Errorf("load data failed for ConfigGetOperation: %+v", err)
		return nil, err
	}
	return config, nil
}

func (r *ConfigRepository) GetAll() ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	return configs, storage.DB.Find(&configs).Error
}
