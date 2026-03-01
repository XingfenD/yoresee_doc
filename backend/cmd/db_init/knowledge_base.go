package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initializeKnowledgeBasesInTx(tx *gorm.DB) error {
	logrus.Println("Initializing default knowledge bases in transaction...")

	var adminUser model.User
	if err := tx.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}

	var count int64
	tx.Model(&model.KnowledgeBase{}).Where("name = ?", "默认知识库").Count(&count)
	if count > 0 {
		logrus.Println("Default knowledge base already exists in transaction.")
		return nil
	}

	// 创建默认知识库
	defaultKnowledgeBase := model.KnowledgeBase{
		ExternalID:    utils.GenerateExternalID("kb"),
		Name:          "默认知识库",
		Description:   "系统默认创建的知识库，用于存储常用文档和资料",
		Cover:         "", // 可以设置默认封面图片URL
		CreatorUserID: adminUser.ID,
		IsPublic:      false,
	}

	if err := tx.Create(&defaultKnowledgeBase).Error; err != nil {
		return err
	}

	// 创建另一个示例知识库
	exampleKnowledgeBase := model.KnowledgeBase{
		ExternalID:    utils.GenerateExternalID("kb"),
		Name:          "示例知识库",
		Description:   "示例知识库，展示知识库功能",
		Cover:         "",
		CreatorUserID: adminUser.ID,
		IsPublic:      true,
	}

	if err := tx.Create(&exampleKnowledgeBase).Error; err != nil {
		return err
	}

	// 尝试将欢迎文档添加到默认知识库中（可选操作，失败不影响整体初始化）
	var welcomeDocument model.DocumentMeta
	if err := tx.Where("title = ?", "欢迎使用 Yoresee Doc").First(&welcomeDocument).Error; err == nil {
		// 如果找到了欢迎文档，则将其关联到默认知识库
		relation := model.DocKnowledgeRelation{
			DocumentID:  welcomeDocument.ID,
			KnowledgeID: &defaultKnowledgeBase.ID,
			OwnerID:     &adminUser.ID,
		}

		if err := tx.Create(&relation).Error; err != nil {
			// 如果关联失败（例如表不存在），记录警告但不中断整个初始化过程
			logrus.Warnf("Could not associate document with knowledge base (may be due to table not existing yet): %v", err)
		} else {
			logrus.Println("Associated welcome document with default knowledge base")
		}
	} else {
		logrus.Println("Welcome document not found, skipping association with knowledge base")
	}

	logrus.Println("Default knowledge bases created successfully in transaction")
	return nil
}
