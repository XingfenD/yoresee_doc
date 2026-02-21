package main

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	registerModeConfigModel := &model.SystemConfig{
		Key: utils.GenConfigKey(
			constant.ConfigKey_First_System,
			constant.ConfigKey_Second_Security,
			constant.ConfigKey_Third_RegisterMode,
		),
		Value: constant.RegisterMode_Open,
	}
	storage.DB.FirstOrCreate(registerModeConfigModel, model.SystemConfig{
		Key: registerModeConfigModel.Key,
	})
}
