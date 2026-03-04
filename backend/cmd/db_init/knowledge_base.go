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

	// 为默认知识库添加11个根文档和子文档
	logrus.Println("Adding documents to default knowledge base...")

	// 根文档标题列表
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

	// 为每个根文档添加子文档
	for _, title := range rootDocumentTitles {
		// 生成根文档的external_id
		rootExternalID := utils.GenerateExternalID("doc")

		// 插入根文档
		var rootDocID int64
		rootDocSQL := `
			INSERT INTO document_metas (external_id, title, type, summary, parent_id, user_id, knowledge_id, status, tags, view_count, edit_count, version, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			RETURNING id
		`
		if err := tx.Raw(rootDocSQL, rootExternalID, title, "markdown", "", 0, adminUser.ID, defaultKnowledgeBase.ID, 1, "[]", 0, 0, 1).Scan(&rootDocID).Error; err != nil {
			logrus.Warnf("Failed to create root document %s: %v", title, err)
			continue
		}

		// 插入文档与知识库的关联关系
		relationSQL := `
			INSERT INTO doc_knowledge_relations (document_id, knowledge_id, owner_id)
			VALUES (?, ?, ?)
		`
		if err := tx.Exec(relationSQL, rootDocID, defaultKnowledgeBase.ID, adminUser.ID).Error; err != nil {
			logrus.Warnf("Failed to create relation for root document %s: %v", title, err)
		}

		// 为每个根文档添加3个子文档
		for j := 1; j <= 3; j++ {
			childTitle := title + " - 子文档" + string(rune('0'+j))
			childExternalID := utils.GenerateExternalID("doc")

			// 插入子文档
			var childDocID int64
			childDocSQL := `
				INSERT INTO document_metas (external_id, title, type, summary, parent_id, user_id, knowledge_id, status, tags, view_count, edit_count, version, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
				RETURNING id
			`
			if err := tx.Raw(childDocSQL, childExternalID, childTitle, "markdown", "", rootDocID, adminUser.ID, defaultKnowledgeBase.ID, 1, "[]", 0, 0, 1).Scan(&childDocID).Error; err != nil {
				logrus.Warnf("Failed to create child document %s: %v", childTitle, err)
				continue
			}

			// 插入子文档与知识库的关联关系
			if err := tx.Exec(relationSQL, childDocID, defaultKnowledgeBase.ID, adminUser.ID).Error; err != nil {
				logrus.Warnf("Failed to create relation for child document %s: %v", childTitle, err)
			}
		}
	}

	logrus.Println("Documents added to default knowledge base successfully")

	logrus.Println("Default knowledge bases created successfully in transaction")
	return nil
}
