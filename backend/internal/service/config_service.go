package service

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

const (
	cacheExpiration   = 5 * time.Minute
	configCachePrefix = "system_config:"
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
	ctx := context.Background()
	cacheKey := configCachePrefix + key

	cachedValue, err := storage.GetCache(ctx, cacheKey)
	if err == nil {
		return cachedValue, nil
	}

	config, err := s.configRepo.Get(key)
	if err != nil {
		return "", err
	}

	storage.SetCache(ctx, cacheKey, config.Value, cacheExpiration)

	return config.Value, nil
}

func (s *ConfigService) Set(key, value string) error {
	var config model.SystemConfig
	if _, err := s.configRepo.Get(key); err != nil {
		config = model.SystemConfig{
			Key:   key,
			Value: value,
		}
		if err := s.configRepo.Create(&config); err != nil {
			return err
		}
	} else {
		config.Value = value

		if err := s.configRepo.Update(&config); err != nil {
			return status.StatusWriteDBError
		}
	}

	ctx := context.Background()
	cacheKey := configCachePrefix + key
	storage.DeleteCache(ctx, cacheKey)

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
