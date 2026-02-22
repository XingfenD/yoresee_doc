package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initializePermissionsInTx(tx *gorm.DB) error {
	logrus.Println("Starting permission initialization in transaction...")

	adminSubject := model.Subject{
		ID:   "admin",
		Type: model.SubjectTypeUser,
	}

	if err := tx.Where("id = ?", adminSubject.ID).FirstOrCreate(&adminSubject).Error; err != nil {
		return err
	}
	logrus.Println("Admin subject created successfully in transaction.")

	logrus.Println("Permission initialization completed successfully in transaction.")
	return nil
}
