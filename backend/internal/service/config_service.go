package service

import (
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
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

func (s *ConfigService) GetSystemRegisterMode() string {
	registerMode, err := s.Get(utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	))
	if err != nil {
		return constant.RegisterMode_Invite
	}
	return registerMode
}

var ConfigSvc = NewConfigService()
