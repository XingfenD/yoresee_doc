package service

import (
	"context"
	"errors"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/lock"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

const (
	cacheExpiration   = 7 * 24 * time.Hour
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

func (s *ConfigService) Get(ctx context.Context, key string) (string, error) {
	cacheKey := configCachePrefix + key

	cachedValue, err := storage.GetCache(ctx, cacheKey)
	if err == nil {
		return cachedValue, nil
	}

	lockKey := lock.LockCachePrefix + cacheKey

	checkFn := func(ctx context.Context) (interface{}, bool, error) {
		cachedValue, err := storage.GetCache(ctx, cacheKey)
		if err == nil {
			return cachedValue, true, nil
		}
		return nil, false, nil
	}

	execFn := func(ctx context.Context) (interface{}, error) {
		cachedValue, err := storage.GetCache(ctx, cacheKey)
		if err == nil {
			return cachedValue, nil
		}

		return s.getConfigAndSetCache(ctx, key, cacheKey)
	}

	result, err := lock.AcquireWithRetry(
		ctx,
		lockKey,
		5*time.Second,
		5,
		100*time.Millisecond,
		checkFn,
		execFn,
	)

	if err != nil {
		return "", err
	}

	if strVal, ok := result.(string); ok {
		return strVal, nil
	}

	return "", status.GenErrWithCustomMsg(status.StatusServiceInternalError, "unexpected return type from lock.AcquireWithRetry")
}

func (s *ConfigService) getConfigAndSetCache(ctx context.Context, key, cacheKey string) (string, error) {
	config, err := s.configRepo.Get(key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			storage.SetCache(ctx, cacheKey, "", 30*time.Second)
			return "", err
		}
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
	storage.SetCache(ctx, cacheKey, value, cacheExpiration)

	return nil
}

func (s *ConfigService) GetSystemRegisterMode(ctx context.Context) string {
	registerMode, err := s.Get(ctx, utils.GenConfigKey(
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
