package main

import (
	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/transport/connectserver"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitRedis().
		InitConsul().
		RequireConsulEnabled().
		InitMinio().
		InitElasticsearchAllowFail().
		InitMQ().
		InitRepository().
		Err(); err != nil {
		logrus.Fatalf("Init backend failed: %v", err)
	}

	defer func() {
		_ = mq.Close()
	}()

	defer func() {
		_ = storage.CloseElasticsearch()
	}()

	defer func() {
		_ = storage.CloseRedis()
	}()

	defer func() {
		_ = storage.ClosePostgres()
	}()

	if err := service.Init(config.GlobalConfig); err != nil {
		logrus.Fatalf("Init services failed: %v", err)
	}

	if _, _, err := connectserver.Start(config.GlobalConfig.Server.GrpcPort, config.GlobalConfig.Server.GrpcWebPort); err != nil {
		logrus.Fatalf("Start connect servers failed: %v", err)
	}

	select {}
}
