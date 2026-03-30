package main

import (
	"context"
	"net/http"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/internal/transport/connectserver"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	initializer := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitRedis().
		InitConsul().
		RequireConsulEnabled().
		InitMinio().
		InitElasticsearchAllowFail().
		InitMQ().
		InitRepository()
	if err := initializer.Err(); err != nil {
		logrus.Fatalf("Init backend failed: %v", err)
	}

	if err := service.Init(config.GlobalConfig); err != nil {
		logrus.Fatalf("Init services failed: %v", err)
	}

	grpcServer, grpcWebServer, err := connectserver.Start(config.GlobalConfig.Server.GrpcPort, config.GlobalConfig.Server.GrpcWebPort)
	if err != nil {
		logrus.Fatalf("Start connect servers failed: %v", err)
	}

	waitForShutdown(initializer, grpcServer, grpcWebServer)
}

func waitForShutdown(initializer *bootstrap.Initializer, grpcServer, grpcWebServer *http.Server) {
	signalCtx, stop := utils.ShutdownContext()
	defer stop()

	<-signalCtx.Done()
	logrus.Info("Shutdown signal received, start draining")
	connectserver.SetDraining(true)

	shutdownServer("grpc-web", grpcWebServer, 10*time.Second)
	shutdownServer("grpc", grpcServer, 10*time.Second)

	initializer.Shutdown()
}

func shutdownServer(name string, server *http.Server, timeout time.Duration) {
	if server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Shutdown %s server failed: %v", name, err)
	}
}
