package config_service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) get(ctx context.Context, key string) (string, bool, error) {
	if !storage.ConsulEnabled() {
		return "", false, errors.New("consul is not enabled")
	}
	return storage.Consul.Get(ctx, key)
}

func (s *ConfigService) GetSystemRegisterMode(ctx context.Context) string {
	registerMode, ok, err := s.get(ctx, utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	))
	if err != nil || !ok || registerMode == "" {
		return constant.RegisterMode_Invite
	}
	return registerMode
}

func (s *ConfigService) GetSystemRegisterLimit(ctx context.Context) bool {
	limitString, ok, err := s.get(ctx, utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterLimit,
	))
	if err != nil || !ok {
		return false
	}
	return limitString == constant.RegisterLimit_On
}

var ConfigSvc = NewConfigService()
