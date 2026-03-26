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

	defaultKnowledgeBase := model.KnowledgeBase{
		ExternalID:    utils.GenerateExternalID("kb"),
		Name:          "默认知识库",
		Description:   "系统默认创建的知识库，用于存储常用文档和资料",
		Cover:         "",
		CreatorUserID: adminUser.ID,
		IsPublic:      false,
	}

	if err := tx.Create(&defaultKnowledgeBase).Error; err != nil {
		return err
	}

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

	logrus.Println("Adding documents to default knowledge base...")

	rootDocumentTitles := []string{
		"项目概述",
		"系统架构",
		"开发规范",
		"API文档",
		"数据库设计",
		"前端开发",
		"后端开发",
		"测试计划",
		"部署指南",
		"安全规范",
		"用户手册",
	}

	for _, title := range rootDocumentTitles {
		rootDoc, err := createDocumentWithAllTables(tx, adminUser.ID, DocumentInput{
			Title:       title,
			Content:     "# " + title + "\n\n这是" + title + "的详细文档内容。",
			Summary:     title + "相关文档",
			ParentID:    0,
			Tags:        []string{title},
			KnowledgeID: &defaultKnowledgeBase.ID,
		})
		if err != nil {
			logrus.Warnf("Failed to create root document %s: %v", title, err)
			continue
		}

		for j := 1; j <= 3; j++ {
			childTitle := title + " - 子文档" + string(rune('0'+j))
			childDoc, err := createDocumentWithAllTables(tx, adminUser.ID, DocumentInput{
				Title:       childTitle,
				Content:     "# " + childTitle + "\n\n这是" + childTitle + "的详细文档内容。",
				Summary:     childTitle + "相关文档",
				ParentID:    rootDoc.ID,
				Tags:        []string{title, "子文档"},
				KnowledgeID: &defaultKnowledgeBase.ID,
			})
			if err != nil {
				logrus.Warnf("Failed to create child document %s: %v", childTitle, err)
				continue
			}

			for k := 1; k <= 2; k++ {
				grandchildTitle := childTitle + " - 孙文档" + string(rune('0'+k))
				_, err := createDocumentWithAllTables(tx, adminUser.ID, DocumentInput{
					Title:       grandchildTitle,
					Content:     "# " + grandchildTitle + "\n\n这是" + grandchildTitle + "的详细文档内容。",
					Summary:     grandchildTitle + "相关文档",
					ParentID:    childDoc.ID,
					Tags:        []string{title, "子文档", "孙文档"},
					KnowledgeID: &defaultKnowledgeBase.ID,
				})
				if err != nil {
					logrus.Warnf("Failed to create grandchild document %s: %v", grandchildTitle, err)
				}
			}
		}
	}

	logrus.Println("Documents added to default knowledge base successfully")

	logrus.Println("Default knowledge bases created successfully in transaction")
	return nil
}
