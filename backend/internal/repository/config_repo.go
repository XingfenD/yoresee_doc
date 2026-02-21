package repository

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type ConfigRepository struct{}

var ConfigRepo = &ConfigRepository{}

func (r *ConfigRepository) Get(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	return &config, storage.DB.Where("key = ?", key).First(&config).Error
}

func (r *ConfigRepository) GetAll() ([]*model.SystemConfig, error) {
	var configs []*model.SystemConfig
	return configs, storage.DB.Find(&configs).Error
}

func (r *ConfigRepository) Create(config *model.SystemConfig) error {
	return storage.DB.Create(config).Error
}

func (r *ConfigRepository) Update(config *model.SystemConfig) error {
	return storage.DB.Save(config).Error
}
