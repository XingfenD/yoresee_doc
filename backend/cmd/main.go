package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/transport/connectserver"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
		panic("init config failed")
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init database failed: %v", err)
		panic("init database failed")
	}

	if err := storage.InitRedis(&config.GlobalConfig.Redis); err != nil {
		logrus.Fatalf("Init redis failed: %v", err)
		panic("init redis failed")
	}

	if err := storage.InitConsul(&config.GlobalConfig.Consul); err != nil {
		logrus.Fatalf("Init consul failed: %v", err)
		panic("init consul failed")
	}
	if !storage.ConsulEnabled() {
		logrus.Fatal("Consul is required for config, but it is not enabled")
	}

	if err := storage.InitMinio(&config.GlobalConfig.Minio); err != nil {
		logrus.Fatalf("Init minio failed: %v", err)
		panic("init minio failed")
	}

	if err := storage.InitElasticsearch(&config.GlobalConfig.Elasticsearch); err != nil {
		logrus.Fatalf("Init elasticsearch failed: %v", err)
		panic("init elasticsearch failed")
	}

	if err := mq.Init(&config.GlobalConfig.MQConfig); err != nil {
		logrus.Fatalf("Init message queue failed: %v", err)
	}

	repository.MustInit()

	defer mq.Close()

	defer storage.ClosePostgres()

	defer storage.CloseRedis()

	defer storage.CloseElasticsearch()

	if err := service.Init(config.GlobalConfig); err != nil {
		logrus.Fatalf("Init services failed: %v", err)
	}

	if _, _, err := connectserver.Start(config.GlobalConfig.Server.GrpcPort, config.GlobalConfig.Server.GrpcWebPort); err != nil {
		logrus.Fatalf("Start connect servers failed: %v", err)
	}

	select {}
}
