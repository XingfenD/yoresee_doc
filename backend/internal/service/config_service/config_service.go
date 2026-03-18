package config_service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/repository/config_repo"
	"github.com/XingfenD/yoresee_doc/internal/utils"
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

func (s *ConfigService) GetSystemRegisterLimit(ctx context.Context) bool {
	limitString, err := s.configRepo.Get((utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterLimit,
	))).Exec(ctx)
	if err != nil {
		return false
	}
	return limitString.Value == constant.RegisterLimit_On
}

var ConfigSvc = NewConfigService()
