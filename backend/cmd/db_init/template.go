package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type templateSeed struct {
	Name            string
	Description     string
	Content         string
	Scope           string
	KnowledgeBaseID *int64
	Tags            []string
}

func initializeTemplatesInTx(tx *gorm.DB) error {
	logrus.Println("Initializing default templates in transaction...")

	var adminUser model.User
	if err := tx.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}

	var defaultKB model.KnowledgeBase
	if err := tx.Where("name = ?", "默认知识库").First(&defaultKB).Error; err != nil {
		return err
	}

	seeds := []templateSeed{
		{
			Name:        "个人周报模板",
			Description: "用于个人每周工作复盘与计划",
			Content: `# 个人周报

## 本周完成
- 

## 问题与风险
- 

## 下周计划
- `,
			Scope: "private",
			Tags:  []string{"个人", "周报"},
		},
		{
			Name:        "会议纪要模板",
			Description: "通用会议记录模板",
			Content: `# 会议纪要

## 会议信息
- 时间：
- 参与人：

## 讨论内容
- 

## 行动项
- [ ] `,
			Scope: "private",
			Tags:  []string{"个人", "会议"},
		},
		{
			Name:            "需求评审模板",
			Description:     "知识库内需求评审记录模板",
			KnowledgeBaseID: &defaultKB.ID,
			Scope:           "knowledge_base",
			Content: `# 需求评审

## 背景

## 目标

## 方案评审
- 优点：
- 风险：

## 结论`,
			Tags: []string{"知识库", "需求"},
		},
		{
			Name:            "技术方案模板",
			Description:     "知识库内技术方案沉淀模板",
			KnowledgeBaseID: &defaultKB.ID,
			Scope:           "knowledge_base",
			Content: `# 技术方案

## 问题定义

## 方案设计

## 数据模型

## 发布计划
`,
			Tags: []string{"知识库", "技术"},
		},
	}

	for _, seed := range seeds {
		if err := createTemplateIfNotExists(tx, adminUser.ID, seed); err != nil {
			return err
		}
	}

	logrus.Println("Default templates initialized successfully in transaction")
	return nil
}

func createTemplateIfNotExists(tx *gorm.DB, userID int64, seed templateSeed) error {
	query := tx.Model(&model.Template{}).
		Where("user_id = ? AND name = ? AND scope = ?", userID, seed.Name, seed.Scope)
	if seed.KnowledgeBaseID == nil {
		query = query.Where("knowledge_base_id IS NULL")
	} else {
		query = query.Where("knowledge_base_id = ?", *seed.KnowledgeBaseID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	template := model.Template{
		Name:            seed.Name,
		Description:     seed.Description,
		DocumentType:    model.DocumentType_Markdown,
		Content:         seed.Content,
		UserID:          userID,
		Scope:           seed.Scope,
		KnowledgeBaseID: seed.KnowledgeBaseID,
		IsPublic:        false,
		Tags:            seed.Tags,
	}
	return tx.Create(&template).Error
}
