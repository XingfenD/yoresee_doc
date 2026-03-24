package main

import (
	"context"
	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

func initializeConfigInConsul(ctx context.Context) error {
	registerModeKey := utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	)

	if _, ok, err := storage.Consul.Get(ctx, registerModeKey); err != nil {
		return err
	} else if !ok {
		if err := storage.Consul.Set(ctx, registerModeKey, constant.RegisterMode_Invite); err != nil {
			return err
		}
	}
	return nil
}

func markDatabaseInitializedInConsul(ctx context.Context) error {
	initializedKey := utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Database,
		constant.ConfigKey_Third_Initialized,
	)
	return storage.Consul.Set(ctx, initializedKey, constant.Database_Initialized_True)
}
