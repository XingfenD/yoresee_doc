package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DocumentInput struct {
	Title       string
	Content     string
	Summary     string
	ParentID    int64
	Tags        []string
	KnowledgeID *int64
}

func createDocumentWithAllTables(tx *gorm.DB, adminUserID int64, input DocumentInput) (*model.DocumentMeta, error) {
	content := model.Content{
		Content: input.Content,
	}

	if err := tx.Create(&content).Error; err != nil {
		return nil, err
	}

	document := model.DocumentMeta{
		ExternalID:  utils.GenerateExternalID("doc"),
		Title:       input.Title,
		Type:        "markdown",
		Summary:     input.Summary,
		ParentID:    input.ParentID,
		UserID:      adminUserID,
		Status:      1,
		Content:     input.Content,
		Tags:        input.Tags,
		ViewCount:   0,
		EditCount:   0,
		KnowledgeID: input.KnowledgeID,
	}

	if err := tx.Create(&document).Error; err != nil {
		return nil, err
	}

	documentVersion := model.DocumentVersion{
		DocumentID:    document.ID,
		Version:       1,
		Title:         document.Title,
		ContentID:     content.ID,
		UserID:        adminUserID,
		ChangeSummary: "Initial version",
	}

	if err := tx.Create(&documentVersion).Error; err != nil {
		return nil, err
	}

	if input.KnowledgeID != nil {
		relation := model.DocKnowledgeRelation{
			DocumentID:  document.ID,
			KnowledgeID: input.KnowledgeID,
			OwnerID:     &adminUserID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			return nil, err
		}
	}

	return &document, nil
}

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
	contentString := "# 欢迎使用 Yoresee Doc\n\n这是您的第一个文档。Yoresee Doc 是一个功能强大的文档管理系统，支持以下特性：\n\n- 📝 富文本编辑\n- 📁 文档分类管理\n- 🔍 全文搜索\n- 👥 协作编辑\n- 📊 版本控制\n- 🔒 权限管理\n\n## 快速开始\n\n1. 点击左侧菜单创建新文档\n2. 使用编辑器撰写内容\n3. 保存文档并分享给团队成员\n\n祝您使用愉快！"
	content := model.Content{
		Content: contentString,
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
		Content:    contentString,
		Tags:       []string{"guide", "welcome"},
		ViewCount:  0,
		EditCount:  0,
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
