package main

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/XingfenD/yoresee_doc/internal/router"
	"github.com/XingfenD/yoresee_doc/internal/service"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init database failed: %v", err)
	}

	if err := storage.InitRedis(&config.GlobalConfig.Redis); err != nil {
		logrus.Fatalf("Init redis failed: %v", err)
	}

	if err := service.InitMessageQueue(config.GlobalConfig); err != nil {
		logrus.Fatalf("Init message queue failed: %v", err)
	}
	defer service.CloseMessageQueue()

	defer storage.ClosePostgres()

	defer storage.CloseRedis()

	i18n.Init()

	r := gin.Default()

	router.SetupRouter(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	r.Run(addr)
}
