package main

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

func main() {
	initializer := bootstrap.NewInitializer().
		InitConfig().
		InitPostgres().
		InitRedis().
		InitMQ().
		InitRepository()
	if err := initializer.Err(); err != nil {
		logrus.Fatalf("Init notification-worker failed: %v", err)
	}

	topic := domain_event.NotificationCreateTopic()

	go func() {
		if err := mq.SubscribeTo(mq.BackendRabbitMQ, topic, handleNotificationEvent); err != nil {
			logrus.Fatalf("Subscribe failed: %v", err)
		}
	}()

	initializer.ShutdownOnSignal(500 * time.Millisecond)
}

func handleNotificationEvent(ctx context.Context, data []byte) error {
	payload := strings.TrimSpace(string(data))
	if payload == "" {
		return nil
	}

	evt, err := domain_event.DecodeNotificationCreateEvent(data)
	if err != nil {
		logrus.Errorf("parse notification event failed: %v", err)
		return nil
	}

	req := &dto.CreateNotificationRequest{
		ReceiverExternalIDs: evt.ReceiverExternalIDs,
		Type:                evt.Type,
		Title:               evt.Title,
		Content:             evt.Content,
		PayloadJSON:         evt.PayloadJSON,
	}
	if err := notification_service.NotificationSvc.CreateNotifications(req); err != nil {
		logrus.Errorf("create notifications failed: %v", err)
	}
	return nil
}
