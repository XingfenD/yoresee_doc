package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initializeDocumentsInTx(tx *gorm.DB) error {
	logrus.Println("Initializing default documents in transaction...")

	var adminUser model.User
	if err := tx.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}

	var count int64
	tx.Model(&model.DocumentMeta{}).Where("title = ?", "欢迎使用 Yoresee Doc").Count(&count)
	if count > 0 {
		logrus.Println("Default document already exists in transaction.")
		return nil
	}

	content := model.Content{
		Content: "# 欢迎使用 Yoresee Doc\n\n这是您的第一个文档。Yoresee Doc 是一个功能强大的文档管理系统，支持以下特性：\n\n- 📝 富文本编辑\n- 📁 文档分类管理\n- 🔍 全文搜索\n- 👥 协作编辑\n- 📊 版本控制\n- 🔒 权限管理\n\n## 快速开始\n\n1. 点击左侧菜单创建新文档\n2. 使用编辑器撰写内容\n3. 保存文档并分享给团队成员\n\n祝您使用愉快！",
	}

	if err := tx.Create(&content).Error; err != nil {
		return err
	}

	document := model.DocumentMeta{
		ExternalID: utils.GenerateExternalID("doc"),
		Title:      "欢迎使用 Yoresee Doc",
		Type:       "markdown",
		Summary:    "Yoresee Doc 系统欢迎文档",
		ParentID:   0,
		UserID:     adminUser.ID,
		Status:     1,
		IsPublic:   true,
		Tags:       []string{"guide", "welcome"},
		ViewCount:  0,
		EditCount:  0,
		Version:    1,
	}

	if err := tx.Create(&document).Error; err != nil {
		return err
	}

	documentVersion := model.DocumentVersion{
		DocumentID:    document.ID,
		Version:       1,
		Title:         document.Title,
		ContentID:     content.ID,
		UserID:        adminUser.ID,
		ChangeSummary: "Initial version",
	}

	if err := tx.Create(&documentVersion).Error; err != nil {
		return err
	}

	logrus.Println("Default document created successfully in transaction")
	return nil
}
