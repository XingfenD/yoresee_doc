package setting_service

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/constant"
	"github.com/XingfenD/yoresee_doc/internal/service/config_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

type SettingOption struct {
	Label    string
	LabelKey string
	Value    string
}

type SettingUI struct {
	Component       string
	Options         []SettingOption
	Placeholder     string
	PlaceholderKey  string
}

type SettingItem struct {
	Key            string
	Label          string
	LabelKey       string
	Description    string
	DescriptionKey string
	Type           string
	UI             SettingUI
	Value          string
	DefaultValue   string
	Required       bool
	Readonly       bool
}

type SettingGroup struct {
	Key      string
	Title    string
	TitleKey string
	Items    []SettingItem
}

type SettingUpdate struct {
	Key   string
	Value string
}

type SettingService struct{}

var SettingSvc *SettingService

func InitSettingService() {
	SettingSvc = &SettingService{}
}

func (s *SettingService) GetSettings(ctx context.Context, scene string) ([]SettingGroup, error) {
	if scene == "" || scene == "system" || scene == "manage" {
		return s.buildSystemSecuritySettings(ctx), nil
	}
	return []SettingGroup{}, nil
}

func (s *SettingService) UpdateSettings(ctx context.Context, updates []SettingUpdate) error {
	for _, update := range updates {
		key := strings.TrimSpace(update.Key)
		if key == "" {
			continue
		}
		switch key {
		case systemRegisterModeKey():
			val := strings.TrimSpace(update.Value)
			if val != constant.RegisterMode_Open && val != constant.RegisterMode_Invite {
				return status.GenErrWithCustomMsg(status.StatusParamError, "invalid register mode")
			}
			if err := config_service.ConfigSvc.Set(ctx, key, val); err != nil {
				logrus.Errorf("[Service layer: SettingService] set register mode failed, key=%s, value=%s, err=%+v", key, val, err)
				return status.GenErrWithCustomMsg(err, "update register mode failed")
			}
			storage.Consul.ClearCacheKey(key)
		case systemRegisterLimitKey():
			val := strings.ToLower(strings.TrimSpace(update.Value))
			boolVal := val == "true" || val == "1" || val == "on"
			writeVal := constant.RegisterLimit_Off
			if boolVal {
				writeVal = constant.RegisterLimit_On
			}
			if err := config_service.ConfigSvc.Set(ctx, key, writeVal); err != nil {
				logrus.Errorf("[Service layer: SettingService] set register limit failed, key=%s, value=%s, err=%+v", key, writeVal, err)
				return status.GenErrWithCustomMsg(err, "update register limit failed")
			}
			storage.Consul.ClearCacheKey(key)
		default:
			return status.GenErrWithCustomMsg(status.StatusParamError, "unknown setting key")
		}
	}
	return nil
}

func (s *SettingService) buildSystemSecuritySettings(ctx context.Context) []SettingGroup {
	registerMode := config_service.ConfigSvc.GetSystemRegisterMode(ctx)
	registerLimit := config_service.ConfigSvc.GetSystemRegisterLimit(ctx)
	limitValue := "false"
	if registerLimit {
		limitValue = "true"
	}

	return []SettingGroup{
		{
			Key:      "security",
			TitleKey: "system.security.registration",
			Items: []SettingItem{
				{
					Key:          systemRegisterModeKey(),
					LabelKey:     "system.security.registrationMode",
					Type:         "string",
					Value:        registerMode,
					DefaultValue: constant.RegisterMode_Invite,
					Required:     true,
					UI: SettingUI{
						Component: "radio",
						Options: []SettingOption{
							{
								LabelKey: "system.security.freeRegister",
								Value:    constant.RegisterMode_Open,
							},
							{
								LabelKey: "system.security.inviteOnly",
								Value:    constant.RegisterMode_Invite,
							},
						},
					},
				},
				{
					Key:          systemRegisterLimitKey(),
					LabelKey:     "system.security.registerLimit",
					DescriptionKey: "system.security.registerLimitDesc",
					Type:         "bool",
					Value:        limitValue,
					DefaultValue: "false",
					UI: SettingUI{
						Component: "switch",
					},
				},
			},
		},
	}
}

func systemRegisterModeKey() string {
	return utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterMode,
	)
}

func systemRegisterLimitKey() string {
	return utils.GenConfigKey(
		constant.ConfigKey_First_System,
		constant.ConfigKey_Second_Security,
		constant.ConfigKey_Third_RegisterLimit,
	)
}
