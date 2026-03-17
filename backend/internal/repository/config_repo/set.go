package config_repo

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func (r *ConfigRepository) Create(config *model.SystemConfig) error {
	// TODO: redis support
	return storage.DB.Create(config).Error
}

func (r *ConfigRepository) Update(config *model.SystemConfig) error {
	// TODO: redis support
	return storage.DB.Save(config).Error
}
