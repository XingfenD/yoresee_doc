package service

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
)

type ConfigService struct {
	configRepo *repository.ConfigRepository
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		configRepo: repository.ConfigRepo,
	}
}

func (s *ConfigService) Get(key string) (string, error) {
	config, err := s.configRepo.Get(key)
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func (s *ConfigService) Set(key, value string) error {
	var config model.SystemConfig
	// check if exists
	if _, err := s.configRepo.Get(key); err != nil {
		// create
		config = model.SystemConfig{
			Key:   key,
			Value: value,
		}
		return s.configRepo.Create(&config)
	}

	// update
	config.Value = value

	if err := s.configRepo.Update(&config); err != nil {
		return status.StatusWriteDBError
	}
	return nil
}

var ConfigSvc = NewConfigService()
