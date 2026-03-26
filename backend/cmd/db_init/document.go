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

func createDocumentWithAllTables(tx *gorm.DB, adminUserID int64, input DocumentInput) (*model.Document, error) {
	document := model.Document{
		ExternalID:  utils.GenerateExternalID(utils.ExternalIDContextDocument),
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
		Path:        "0",
		Depth:       0,
	}

	if err := tx.Create(&document).Error; err != nil {
		return nil, err
	}
	if err := updateDocumentPathDepth(tx, document.ID, document.ParentID); err != nil {
		return nil, err
	}

	documentVersion := model.DocumentVersion{
		DocumentID:    document.ID,
		Version:       1,
		Title:         document.Title,
		Content:       document.Content,
		UserID:        adminUserID,
		ChangeSummary: "Initial version",
	}

	if err := tx.Create(&documentVersion).Error; err != nil {
		return nil, err
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
	tx.Model(&model.Document{}).Where("title = ?", "欢迎使用 Yoresee Doc").Count(&count)
	if count > 0 {
		logrus.Println("Default document already exists in transaction.")
		return nil
	}
	contentString := "# 欢迎使用 Yoresee Doc\n\n这是您的第一个文档。Yoresee Doc 是一个功能强大的文档管理系统，支持以下特性：\n\n- 📝 富文本编辑\n- 📁 文档分类管理\n- 🔍 全文搜索\n- 👥 协作编辑\n- 📊 版本控制\n- 🔒 权限管理\n\n## 快速开始\n\n1. 点击左侧菜单创建新文档\n2. 使用编辑器撰写内容\n3. 保存文档并分享给团队成员\n\n祝您使用愉快！"

	document := model.Document{
		ExternalID: utils.GenerateExternalID(utils.ExternalIDContextDocument),
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
		Path:       "0",
		Depth:      0,
	}

	if err := tx.Create(&document).Error; err != nil {
		return err
	}
	if err := updateDocumentPathDepth(tx, document.ID, document.ParentID); err != nil {
		return err
	}

	documentVersion := model.DocumentVersion{
		DocumentID:    document.ID,
		Version:       1,
		Title:         document.Title,
		Content:       contentString,
		UserID:        adminUser.ID,
		ChangeSummary: "Initial version",
	}

	if err := tx.Create(&documentVersion).Error; err != nil {
		return err
	}

	logrus.Println("Adding child documents to welcome document...")

	childDocumentTitles := []string{
		"功能介绍",
		"使用教程",
		"常见问题",
		"更新日志",
		"联系我们",
	}

	for _, title := range childDocumentTitles {
		childDoc, err := createDocumentWithAllTables(tx, adminUser.ID, DocumentInput{
			Title:       title,
			Content:     "# " + title + "\n\n这是" + title + "的详细文档内容。",
			Summary:     title + "相关文档",
			ParentID:    document.ID,
			Tags:        []string{"welcome", title},
			KnowledgeID: nil,
		})
		if err != nil {
			logrus.Warnf("Failed to create child document %s: %v", title, err)
			continue
		}

		for k := 1; k <= 2; k++ {
			grandchildTitle := title + " - 详情" + string(rune('0'+k))
			_, err := createDocumentWithAllTables(tx, adminUser.ID, DocumentInput{
				Title:       grandchildTitle,
				Content:     "# " + grandchildTitle + "\n\n这是" + grandchildTitle + "的详细文档内容。",
				Summary:     grandchildTitle + "相关文档",
				ParentID:    childDoc.ID,
				Tags:        []string{"welcome", title, "详情"},
				KnowledgeID: nil,
			})
			if err != nil {
				logrus.Warnf("Failed to create grandchild document %s: %v", grandchildTitle, err)
			}
		}
	}

	logrus.Println("Child documents added to welcome document successfully")

	logrus.Println("Default document created successfully in transaction")
	return nil
}

func updateDocumentPathDepth(tx *gorm.DB, docID int64, parentID int64) error {
	if parentID == 0 {
		return tx.Exec(`
			UPDATE document_metas
			SET path = (id::text)::ltree, depth = 0
			WHERE id = ?
		`, docID).Error
	}

	type pathDepth struct {
		Path  string
		Depth int
	}
	var parent pathDepth
	if err := tx.Model(&model.Document{}).
		Select("path, depth").
		Where("id = ?", parentID).
		Take(&parent).Error; err != nil {
		return err
	}

	return tx.Exec(`
		UPDATE document_metas
		SET path = (?::ltree) || (id::text)::ltree, depth = ?
		WHERE id = ?
	`, parent.Path, parent.Depth+1, docID).Error
}
