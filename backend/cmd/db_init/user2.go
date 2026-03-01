package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func createUser2InTx(tx *gorm.DB) error {
	logrus.Println("Creating user2 in transaction...")

	password := "user2pass123"
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	externalID := utils.GenerateExternalID("usr")

	user2 := model.User{
		ExternalID:   externalID,
		Username:     "user2",
		PasswordHash: hashedPwd,
		Email:        "user2@example.com",
		Nickname:     "User Two",
		Status:       1,
	}

	var count int64
	tx.Model(&model.User{}).Where("email = ?", user2.Email).Count(&count)
	if count == 0 {
		if err := tx.Create(&user2).Error; err != nil {
			return err
		}
		logrus.Printf("User2 created successfully with ID: %s in transaction.\n", user2.ExternalID)
	} else {
		logrus.Println("User2 already exists in transaction.")
		// 获取已存在的用户信息
		if err := tx.Where("email = ?", user2.Email).First(&user2).Error; err != nil {
			return err
		}
	}

	// 使用用户2的身份创建一个知识库
	logrus.Println("Creating knowledge base for user2...")

	kbExternalID := utils.GenerateExternalID("kb")
	knowledgeBase := model.KnowledgeBase{
		ExternalID:    kbExternalID,
		Name:          "User2's Knowledge Base",
		Description:   "Knowledge base created by user2",
		CreatorUserID: user2.ID, // 使用用户2的ID
		IsPublic:      false,
	}

	// 检查知识库是否已存在
	var kbCount int64
	tx.Model(&model.KnowledgeBase{}).Where("name = ?", knowledgeBase.Name).Count(&kbCount)
	if kbCount == 0 {
		if err := tx.Create(&knowledgeBase).Error; err != nil {
			return err
		}
		logrus.Printf("Knowledge base created successfully for user2: %s\n", knowledgeBase.Name)
	} else {
		logrus.Println("Knowledge base for user2 already exists.")
	}

	return nil
}
