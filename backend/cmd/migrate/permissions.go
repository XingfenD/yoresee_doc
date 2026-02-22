package main

import (
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func initializePermissions() error {
	logrus.Println("Starting permission initialization...")

	defaultNamespace := model.Namespace{
		ID:      "default",
		Name:    "Default Namespace",
		OwnerID: "admin",
	}

	if err := storage.DB.Where("id = ?", defaultNamespace.ID).FirstOrCreate(&defaultNamespace).Error; err != nil {
		return err
	}
	logrus.Println("Default namespace created successfully.")

	adminSubject := model.Subject{
		ID:        "admin",
		Namespace: "default",
		Type:      model.SubjectTypeUser,
	}

	if err := storage.DB.Where("id = ?", adminSubject.ID).FirstOrCreate(&adminSubject).Error; err != nil {
		return err
	}
	logrus.Println("Admin subject created successfully.")

	namespaceResource := model.Resource{
		ID:        "default",
		Namespace: "default",
		Type:      model.ResourceTypeNamespace,
	}

	if err := storage.DB.Where("id = ?", namespaceResource.ID).FirstOrCreate(&namespaceResource).Error; err != nil {
		return err
	}
	logrus.Println("Namespace resource created successfully.")

	logrus.Println("Permission initialization completed successfully.")
	return nil
}
