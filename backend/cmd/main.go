package main

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/XingfenD/yoresee_doc/internal/router"
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
