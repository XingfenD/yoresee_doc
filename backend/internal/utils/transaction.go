package utils

import (
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

func WithTransaction(fn func(tx *gorm.DB) error) error {
	tx := storage.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
