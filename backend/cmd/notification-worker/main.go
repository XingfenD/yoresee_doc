package main

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/bootstrap"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
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
	group := utils.GetEnvVar("NOTIFICATION_MQ_GROUP", "notification-worker")

	go func() {
		if err := mq_service.MQSvc.Consume(
			context.Background(),
			mq.BackendRabbitMQ,
			mq.ConsumeOptions{
				Topic:   topic,
				Mode:    mq.ConsumeModeGroup,
				Group:   group,
				AutoAck: false,
				OnError: mq.ErrorActionRequeue,
			},
			func(ctx context.Context, message mq.Message) error {
				return handleNotificationEvent(ctx, message.Body)
			},
		); err != nil {
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
		return err
	}
	return nil
}
