package main

import (
	"strconv"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createTestUserIntx(tx *gorm.DB) error {
	logrus.Println("Creating test user in transaction...")
	return createUserWithUsername(tx, "test")
}

func createNormalUserInTx(tx *gorm.DB) error {
	logrus.Println("Creating 100 users in transaction...")

	for i := 1; i <= 20; i++ {
		username := "user" + strconv.Itoa(i)
		if err := createUserWithUsername(tx, username); err != nil {
			return err
		}

		var user model.User
		if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
			return err
		}
		userNum := strconv.Itoa(i)
		logrus.Printf("User%s created successfully with ID: %s in transaction.\n", userNum, user.ExternalID)

		logrus.Printf("Creating knowledge base for user%s...\n", userNum)

		kbExternalID := utils.GenerateExternalID(utils.ExternalIDKnowledgeBase)
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
