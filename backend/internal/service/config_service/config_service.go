package config_service

import (
	"context"
	"sync"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type ConfigService struct {
	SystemRegisterMode  func() string `consul:"system.security.register_mode,default=invite"`
	SystemRegisterLimit func() bool   `consul:"system.security.register_limit"`
	boundOnce           sync.Once
}

func NewConfigService() *ConfigService {
	s := &ConfigService{}
	if storage.ConsulEnabled() {
		s.boundOnce.Do(func() {
			_ = storage.BindConsulConfig(s, storage.Consul)
		})
	}
	return s
}

func (s *ConfigService) GetSystemRegisterMode(ctx context.Context) string {
	return s.SystemRegisterMode()
}

func (s *ConfigService) GetSystemRegisterLimit(ctx context.Context) bool {
	return s.SystemRegisterLimit()
}

var ConfigSvc *ConfigService

func InitConfigService() {
	ConfigSvc = NewConfigService()
}
