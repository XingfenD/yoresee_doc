package config_service

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type ConfigService struct {
	SystemRegisterMode  func() string `consul:"system.security.register_mode,default=invite"`
	SystemRegisterLimit func() bool   `consul:"system.security.register_limit"`
}

func NewConfigService() *ConfigService {
	if !storage.ConsulEnabled() {
		panic("consul is required for config")
	}
	s := &ConfigService{}
	if err := storage.BindConsulConfig(s, storage.Consul); err != nil {
		panic("consul service init failed")
	}
	return s
}

func (s *ConfigService) GetSystemRegisterMode(ctx context.Context) string {
	return s.SystemRegisterMode()
}

func (s *ConfigService) GetSystemRegisterLimit(ctx context.Context) bool {
	return s.SystemRegisterLimit()
}

func (s *ConfigService) Set(ctx context.Context, key, value string) error {
	if !storage.ConsulEnabled() {
		return errors.New("consul is not enabled")
	}
	return storage.Consul.Set(ctx, key, value)
}

var ConfigSvc *ConfigService

func InitConfigService() {
	ConfigSvc = NewConfigService()
}
