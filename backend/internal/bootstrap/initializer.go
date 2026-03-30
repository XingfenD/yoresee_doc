package bootstrap

import (
	"fmt"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
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
	i.Check("storage.InitPostgres", func() error {
		return storage.InitPostgres(&config.GlobalConfig.Database)
	})
	return i.addShutdownHook("Postgres", storage.ClosePostgres)
}

func (i *Initializer) InitRedis() *Initializer {
	i.Check("storage.InitRedis", func() error {
		return storage.InitRedis(&config.GlobalConfig.Redis)
	})
	return i.addShutdownHook("Redis", storage.CloseRedis)
}

func (i *Initializer) InitConsul() *Initializer {
	return i.Check("storage.InitConsul", func() error {
		return storage.InitConsul(&config.GlobalConfig.Consul)
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
	return i.Check("storage.InitMinio", func() error {
		return storage.InitMinio(&config.GlobalConfig.Minio)
	})
}

func (i *Initializer) InitElasticsearch() *Initializer {
	i.Check("storage.InitElasticsearch", func() error {
		return storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch)
	})
	return i.addShutdownHook("Elasticsearch", storage.CloseElasticsearch)
}

func (i *Initializer) InitElasticsearchAllowFail() *Initializer {
	i.CheckAllowFail("storage.InitElasticsearch", func() error {
		return storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch)
	})
	if storage.ES != nil {
		i.addShutdownHook("Elasticsearch", storage.CloseElasticsearch)
	}
	return i
}

func (i *Initializer) InitMQ() *Initializer {
	i.Check("mq.Init", func() error {
		return mq.Init(&config.GlobalConfig.MQConfig)
	})
	return i.addShutdownHook("MQ", mq.Close)
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
