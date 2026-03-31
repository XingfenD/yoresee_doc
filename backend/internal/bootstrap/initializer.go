package bootstrap

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

type Initializer struct {
	err           error
	shutdownHooks []utils.ShutdownHook
	hookNames     map[string]struct{}
	shutdownOnce  sync.Once
}

func NewInitializer() *Initializer {
	return &Initializer{
		hookNames: make(map[string]struct{}),
	}
}

func (i *Initializer) Check(step string, fn func() error) *Initializer {
	if i.err != nil {
		return i
	}
	if err := fn(); err != nil {
		i.err = fmt.Errorf("%s failed: %w", step, err)
	}
	return i
}

func (i *Initializer) CheckAllowFail(step string, fn func() error) *Initializer {
	if i.err != nil {
		return i
	}
	if err := fn(); err != nil {
		logrus.Warnf("%s failed, continue with degraded mode: %v", step, err)
	}
	return i
}

func (i *Initializer) InitConfig() *Initializer {
	return i.Check("config.InitConfig", config.InitConfig)
}

func (i *Initializer) InitPostgres() *Initializer {
	i.Check("storage.OpenPostgres", func() error {
		dbCfg := config.GlobalConfig.Database
		logCfg := config.GlobalConfig.Backend.Log
		gormLogLevel := logCfg.GormLogLevel
		if gormLogLevel == "" {
			gormLogLevel = logCfg.Level
		}

		db, err := storage.OpenPostgres(&storage.PostgresOptions{
			Host:         dbCfg.Host,
			Port:         dbCfg.Port,
			User:         dbCfg.User,
			Password:     dbCfg.Password,
			Name:         dbCfg.Name,
			MaxIdleConns: dbCfg.MaxIdleConns,
			MaxOpenConns: dbCfg.MaxOpenConns,
			GormLogLevel: gormLogLevel,
		})
		if err != nil {
			return err
		}
		storage.DB = db
		return nil
	})

	return i.addShutdownHook("Postgres", func() error {
		if storage.DB == nil {
			return nil
		}
		sqlDB, err := storage.DB.DB()
		if err != nil {
			return err
		}
		storage.DB = nil
		return sqlDB.Close()
	})
}

func (i *Initializer) InitRedis() *Initializer {
	i.Check("storage.NewRedisClient", func() error {
		cfg := config.GlobalConfig.Redis
		client, err := storage.NewRedisClient(&storage.RedisOptions{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
		if err != nil {
			return err
		}
		storage.KVS = client
		return nil
	})

	return i.addShutdownHook("Redis", func() error {
		if storage.KVS == nil {
			return nil
		}
		client := storage.KVS
		storage.KVS = nil
		return client.Close()
	})
}

func (i *Initializer) InitConsul() *Initializer {
	return i.Check("storage.NewConsulClient", func() error {
		cfg := config.GlobalConfig.Consul
		storage.Consul = storage.NewConsulClient(&storage.ConsulOptions{
			Enabled:    cfg.Enabled,
			Address:    cfg.Address,
			Scheme:     cfg.Scheme,
			Token:      cfg.Token,
			Datacenter: cfg.Datacenter,
			Prefix:     cfg.Prefix,
		})
		return nil
	})
}

func (i *Initializer) RequireConsulEnabled() *Initializer {
	return i.Check("storage.ConsulEnabled", func() error {
		if !storage.ConsulEnabled() {
			return fmt.Errorf("consul is required but not enabled")
		}
		return nil
	})
}

func (i *Initializer) InitMinio() *Initializer {
	return i.Check("storage.NewMinioClient", func() error {
		cfg := config.GlobalConfig.Minio
		client, err := storage.NewMinioClient(&storage.MinioOptions{
			Endpoint:  cfg.Endpoint,
			AccessKey: cfg.AccessKey,
			SecretKey: cfg.SecretKey,
			UseSSL:    cfg.UseSSL,
		})
		if err != nil {
			return err
		}
		if err = storage.EnsureBucket(context.Background(), client, cfg.Bucket); err != nil {
			return err
		}
		storage.MinioClient = client
		return nil
	})
}

func (i *Initializer) InitElasticsearch() *Initializer {
	i.Check("storage.NewElasticsearchClient", func() error {
		cfg := config.GlobalConfig.Elasticsearch
		client, err := storage.NewElasticsearchClient(&storage.ElasticsearchOptions{
			Enabled:   cfg.Enabled,
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
			Timeout:   cfg.Timeout,
		})
		if err != nil {
			return err
		}
		storage.ES = client
		return nil
	})

	return i.addShutdownHook("Elasticsearch", func() error {
		storage.ES = nil
		return nil
	})
}

func (i *Initializer) InitElasticsearchAllowFail() *Initializer {
	i.CheckAllowFail("storage.NewElasticsearchClient", func() error {
		cfg := config.GlobalConfig.Elasticsearch
		client, err := storage.NewElasticsearchClient(&storage.ElasticsearchOptions{
			Enabled:   cfg.Enabled,
			Addresses: cfg.Addresses,
			Username:  cfg.Username,
			Password:  cfg.Password,
			Timeout:   cfg.Timeout,
		})
		if err != nil {
			return err
		}
		storage.ES = client
		return nil
	})
	if storage.ES != nil {
		i.addShutdownHook("Elasticsearch", func() error {
			storage.ES = nil
			return nil
		})
	}
	return i
}

func (i *Initializer) InitMQ() *Initializer {
	i.Check("mq_service.Init", func() error {
		return mq_service.MQSvc.Init(&config.GlobalConfig.MQConfig)
	})
	return i.addShutdownHook("MQ", mq_service.MQSvc.Close)
}

func (i *Initializer) InitRepository() *Initializer {
	return i.Check("repository.MustInit", func() error {
		repository.MustInit()
		return nil
	})
}

func (i *Initializer) Err() error {
	return i.err
}

func (i *Initializer) Shutdown() {
	i.shutdownOnce.Do(func() {
		if len(i.shutdownHooks) == 0 {
			return
		}
		reversed := make([]utils.ShutdownHook, 0, len(i.shutdownHooks))
		for idx := len(i.shutdownHooks) - 1; idx >= 0; idx-- {
			reversed = append(reversed, i.shutdownHooks[idx])
		}
		utils.RunShutdownHooks(reversed...)
	})
}

func (i *Initializer) ShutdownOnSignal(delay time.Duration) {
	utils.WaitForShutdownSignal()
	if delay > 0 {
		time.Sleep(delay)
	}
	i.Shutdown()
}

func (i *Initializer) addShutdownHook(name string, fn func() error) *Initializer {
	if i.err != nil || fn == nil || name == "" {
		return i
	}
	if _, exists := i.hookNames[name]; exists {
		return i
	}
	i.hookNames[name] = struct{}{}
	i.shutdownHooks = append(i.shutdownHooks, utils.ShutdownHook{
		Name: name,
		Fn:   fn,
	})
	return i
}
