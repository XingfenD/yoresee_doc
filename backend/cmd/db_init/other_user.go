package main

import (
	"strconv"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createNormalUserInTx(tx *gorm.DB) error {
	logrus.Println("Creating 100 users in transaction...")

	password := "user2pass123"
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	for i := 1; i <= 100; i++ {
		userNum := strconv.Itoa(i)
		externalID := utils.GenerateExternalID("usr")

		user := model.User{
			ExternalID:   externalID,
			Username:     "user" + userNum,
			PasswordHash: hashedPwd,
			Email:        "user" + userNum + "@yoresee.cc",
			Nickname:     "User " + userNum,
			Status:       1,
		}

		var count int64
		tx.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)
		if count == 0 {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			logrus.Printf("User%s created successfully with ID: %s in transaction.\n", userNum, user.ExternalID)
		} else {
			logrus.Printf("User%s already exists in transaction.\n", userNum)
			if err := tx.Where("email = ?", user.Email).First(&user).Error; err != nil {
				return err
			}
		}

		logrus.Printf("Creating knowledge base for user%s...\n", userNum)

		kbExternalID := utils.GenerateExternalID("kb")
		knowledgeBase := model.KnowledgeBase{
			ExternalID:    kbExternalID,
			Name:          "User" + userNum + "'s Knowledge Base",
			Description:   "Knowledge base created by user" + userNum,
			CreatorUserID: user.ID,
			IsPublic:      false,
		}

		var kbCount int64
		tx.Model(&model.KnowledgeBase{}).Where("name = ?", knowledgeBase.Name).Count(&kbCount)
		if kbCount == 0 {
			if err := tx.Create(&knowledgeBase).Error; err != nil {
				return err
			}
			logrus.Printf("Knowledge base created successfully for user%s: %s\n", userNum, knowledgeBase.Name)
		} else {
			logrus.Printf("Knowledge base for user%s already exists.\n", userNum)
		}
	}

	return nil
}
