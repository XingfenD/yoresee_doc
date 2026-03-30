package bootstrap

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

type Initializer struct {
	err error
}

func NewInitializer() *Initializer {
	return &Initializer{}
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
	return i.Check("storage.InitPostgres", func() error {
		return storage.InitPostgres(&config.GlobalConfig.Database)
	})
}

func (i *Initializer) InitRedis() *Initializer {
	return i.Check("storage.InitRedis", func() error {
		return storage.InitRedis(&config.GlobalConfig.Redis)
	})
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
	return i.Check("storage.InitElasticsearch", func() error {
		return storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch)
	})
}

func (i *Initializer) InitElasticsearchAllowFail() *Initializer {
	return i.CheckAllowFail("storage.InitElasticsearch", func() error {
		return storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch)
	})
}

func (i *Initializer) InitMQ() *Initializer {
	return i.Check("mq.Init", func() error {
		return mq.Init(&config.GlobalConfig.MQConfig)
	})
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
