package main

import (
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initializeConfigInTx(tx *gorm.DB) error {
	logrus.Println("Initializing system config in transaction...")

	registerModeConfigModel := &model.SystemConfig{
		Key: utils.GenConfigKey(
			constant.ConfigKey_First_System,
			constant.ConfigKey_Second_Security,
			constant.ConfigKey_Third_RegisterMode,
		),
		Value: constant.RegisterMode_Open,
	}

	if err := tx.FirstOrCreate(registerModeConfigModel, model.SystemConfig{
		Key: registerModeConfigModel.Key,
	}).Error; err != nil {
		return err
	}

	logrus.Println("System config initialized successfully in transaction")
	return nil
}
