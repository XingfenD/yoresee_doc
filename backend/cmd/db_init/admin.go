package main

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createUserInTx(tx *gorm.DB, username, email, password string) error {
	logrus.Printf("Creating user %s in transaction...", username)

	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	externalID := utils.GenerateExternalID(utils.ExternalIDContextDocument)

	user := model.User{
		ExternalID:   externalID,
		Username:     username,
		PasswordHash: hashedPwd,
		Email:        email,
		Nickname:     username,
		Status:       1,
	}

	var count int64
	tx.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
	if count == 0 {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		logrus.Printf("User %s created successfully in transaction.", username)
	} else {
		logrus.Printf("User %s already exists in transaction.", username)
	}

	return nil
}

func createUserWithUsername(tx *gorm.DB, username string) error {
	email := fmt.Sprintf("%s@yoresee.cc", username)
	password := username
	return createUserInTx(tx, username, email, password)
}

func createAdminUserInTx(tx *gorm.DB) error {
	logrus.Println("Creating admin user in transaction...")

	if err := createUserWithUsername(tx, "admin"); err != nil {
		return err
	}

	logrus.Println("Admin permission granted successfully in transaction.")
	return nil
}
