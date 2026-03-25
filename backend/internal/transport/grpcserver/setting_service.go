package grpcserver

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/setting_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
)

type SettingServiceServer struct {
	pb.UnimplementedSettingServiceServer
}

func NewSettingServiceServer() *SettingServiceServer {
	return &SettingServiceServer{}
}

func (s *SettingServiceServer) requireAdmin(ctx context.Context) error {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || userExternalID == "" {
		return status.StatusTokenInvalid
	}
	isAdmin, err := auth_service.AuthSvc.IsAdmin(userExternalID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return status.StatusPermissionDenied
	}
	return nil
}

func (s *SettingServiceServer) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {
	if req == nil {
		return &pb.GetSettingsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.GetSettingsResponse{Base: baseResponseFromErr(err)}, nil
	}

	scene := strings.TrimSpace(req.Scene)
	groups, err := setting_service.SettingSvc.GetSettings(ctx, scene)
	if err != nil {
		return &pb.GetSettingsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respGroups := make([]*pb.SettingGroup, 0, len(groups))
	for _, group := range groups {
		items := make([]*pb.SettingItem, 0, len(group.Items))
		for _, item := range group.Items {
			options := make([]*pb.SettingOption, 0, len(item.UI.Options))
			for _, opt := range item.UI.Options {
				options = append(options, &pb.SettingOption{
					Label:    opt.Label,
					LabelKey: opt.LabelKey,
					Value:    opt.Value,
				})
			}
			items = append(items, &pb.SettingItem{
				Key:            item.Key,
				Label:          item.Label,
				LabelKey:       item.LabelKey,
				Description:    item.Description,
				DescriptionKey: item.DescriptionKey,
				Type:           item.Type,
				Ui: &pb.SettingUI{
					Component:      item.UI.Component,
					Options:        options,
					Placeholder:    item.UI.Placeholder,
					PlaceholderKey: item.UI.PlaceholderKey,
				},
				Value:        item.Value,
				DefaultValue: item.DefaultValue,
				Required:     item.Required,
				Readonly:     item.Readonly,
			})
		}
		respGroups = append(respGroups, &pb.SettingGroup{
			Key:      group.Key,
			Title:    group.Title,
			TitleKey: group.TitleKey,
			Items:    items,
		})
	}

	return &pb.GetSettingsResponse{
		Base:   baseResponseFromErr(nil),
		Groups: respGroups,
	}, nil
}

func (s *SettingServiceServer) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.UpdateSettingsResponse, error) {
	if req == nil {
		return &pb.UpdateSettingsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	if err := s.requireAdmin(ctx); err != nil {
		return &pb.UpdateSettingsResponse{Base: baseResponseFromErr(err)}, nil
	}

	updates := make([]setting_service.SettingUpdate, 0, len(req.Updates))
	for _, update := range req.Updates {
		if update == nil {
			continue
		}
		updates = append(updates, setting_service.SettingUpdate{
			Key:   update.Key,
			Value: update.Value,
		})
	}
	if err := setting_service.SettingSvc.UpdateSettings(ctx, updates); err != nil {
		return &pb.UpdateSettingsResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.UpdateSettingsResponse{Base: baseResponseFromErr(nil)}, nil
}
