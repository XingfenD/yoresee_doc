package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Addr            string `env:"ADDR"              envDefault:":1234"`
	JWTSecret       string `env:"JWT_SECRET"`
	CollabCoreURL   string `env:"COLLAB_CORE_URL"   envDefault:"ws://collab-core:1234"`
	BackendGRPCAddr string `env:"BACKEND_GRPC_ADDR" envDefault:"backend:9090"`
	InternalRPCKey  string `env:"INTERNAL_RPC_KEY"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
