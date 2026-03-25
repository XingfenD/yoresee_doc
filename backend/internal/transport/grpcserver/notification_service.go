package grpcserver

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/internal/service/auth_service"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/constant"
	pb "github.com/XingfenD/yoresee_doc/pkg/gen/yoresee_doc/v1"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationServiceServer() *NotificationServiceServer {
	return &NotificationServiceServer{}
}

func (s *NotificationServiceServer) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	if req == nil {
		return &pb.CreateNotificationResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}

	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.CreateNotificationResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	receiverExternalIDs := req.ReceiverExternalIds
	if len(receiverExternalIDs) == 0 {
		receiverExternalIDs = []string{userExternalID}
	}

	if !onlySelf(receiverExternalIDs, userExternalID) {
		isAdmin, err := auth_service.AuthSvc.IsAdmin(userExternalID)
		if err != nil {
			return &pb.CreateNotificationResponse{Base: baseResponseFromErr(err)}, nil
		}
		if !isAdmin {
			return &pb.CreateNotificationResponse{Base: baseResponseFromStatus(status.StatusPermissionDenied)}, nil
		}
	}

	evt := notificationEvent{
		ReceiverExternalIDs: receiverExternalIDs,
		Type:                req.Type,
		Title:               req.Title,
		Content:             req.Content,
		PayloadJSON:         req.PayloadJson,
	}
	data, err := json.Marshal(evt)
	if err != nil {
		return &pb.CreateNotificationResponse{Base: baseResponseFromErr(status.StatusInternalParamsError)}, nil
	}
	topic := constant.NotificationTopicDefault

	if err := mq.PublishTo(ctx, mq.BackendRabbitMQ, topic, data); err != nil {
		return &pb.CreateNotificationResponse{Base: baseResponseFromErr(status.StatusMQNotInitialized)}, nil
	}

	return &pb.CreateNotificationResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *NotificationServiceServer) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	if req == nil {
		return &pb.ListNotificationsResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.ListNotificationsResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	var statusFilter *string
	if req.Status != nil && strings.TrimSpace(req.GetStatus()) != "" {
		statusFilter = req.Status
	}

	list, total, err := notification_service.NotificationSvc.ListNotifications(&dto.ListNotificationsRequest{
		UserExternalID: userExternalID,
		Status:         statusFilter,
		Pagination: dto.Pagination{
			Page:     int(req.Page),
			PageSize: int(req.PageSize),
		},
	})
	if err != nil {
		return &pb.ListNotificationsResponse{Base: baseResponseFromErr(err)}, nil
	}

	respItems := make([]*pb.Notification, 0, len(list))
	for _, item := range list {
		respItems = append(respItems, toNotificationResponse(&item))
	}

	return &pb.ListNotificationsResponse{
		Base:          baseResponseFromErr(nil),
		Notifications: respItems,
		Total:         total,
	}, nil
}

func (s *NotificationServiceServer) MarkNotificationsRead(ctx context.Context, req *pb.MarkNotificationsReadRequest) (*pb.MarkNotificationsReadResponse, error) {
	if req == nil {
		return &pb.MarkNotificationsReadResponse{Base: baseResponseFromStatus(status.StatusParamError)}, nil
	}
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.MarkNotificationsReadResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}

	err := notification_service.NotificationSvc.MarkRead(&dto.MarkNotificationsReadRequest{
		UserExternalID: userExternalID,
		IDs:            req.Ids,
	})
	if err != nil {
		return &pb.MarkNotificationsReadResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.MarkNotificationsReadResponse{Base: baseResponseFromErr(nil)}, nil
}

func (s *NotificationServiceServer) MarkAllNotificationsRead(ctx context.Context, req *pb.MarkAllNotificationsReadRequest) (*pb.MarkAllNotificationsReadResponse, error) {
	userExternalID, ok := ctx.Value("user_external_id").(string)
	if !ok || strings.TrimSpace(userExternalID) == "" {
		return &pb.MarkAllNotificationsReadResponse{Base: baseResponseFromStatus(status.StatusTokenInvalid)}, nil
	}
	if err := notification_service.NotificationSvc.MarkAllRead(userExternalID); err != nil {
		return &pb.MarkAllNotificationsReadResponse{Base: baseResponseFromErr(err)}, nil
	}
	return &pb.MarkAllNotificationsReadResponse{Base: baseResponseFromErr(nil)}, nil
}

func toNotificationResponse(item *model.Notification) *pb.Notification {
	if item == nil {
		return nil
	}
	return &pb.Notification{
		Id:        item.ID,
		Type:      item.Type,
		Status:    item.Status,
		Title:     item.Title,
		Content:   item.Content,
		Payload:   item.Payload,
		CreatedAt: timeToString(item.CreatedAt),
	}
}

func onlySelf(list []string, self string) bool {
	if len(list) == 0 {
		return true
	}
	for _, id := range list {
		if strings.TrimSpace(id) == "" {
			continue
		}
		if id != self {
			return false
		}
	}
	return true
}

type notificationEvent struct {
	ReceiverExternalIDs []string `json:"receiver_external_ids"`
	Type                string   `json:"type"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	PayloadJSON         string   `json:"payload_json"`
}
