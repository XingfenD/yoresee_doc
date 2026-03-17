package config_service

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/repository/config_repo"
	"github.com/XingfenD/yoresee_doc/internal/utils"
)

const (
	cacheExpiration   = 7 * 24 * time.Hour
	configCachePrefix = "system_config:"
)

type ConfigService struct {
	configRepo *config_repo.ConfigRepository
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		configRepo: &config_repo.ConfigRepo,
	}
}
func (s *ConfigService) GetSystemRegisterMode(ctx context.Context) string {
	registerMode, err := s.configRepo.Get(utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	)).Exec(ctx)
	if err != nil {
		return constant.RegisterMode_Invite
	}
	return registerMode.Value
}

var ConfigSvc = NewConfigService()
