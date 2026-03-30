package main

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/domain_event"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig(); err != nil {
		logrus.Fatalf("Init config failed: %v", err)
	}

	if err := storage.InitPostgres(&config.GlobalConfig.Database); err != nil {
		logrus.Fatalf("Init Postgres failed: %v", err)
	}

	if err := storage.InitRedis(&config.GlobalConfig.Redis); err != nil {
		logrus.Fatalf("Init Redis failed: %v", err)
	}

	if err := mq.Init(&config.GlobalConfig.MQConfig); err != nil {
		logrus.Fatalf("Init MQ failed: %v", err)
	}

	repository.MustInit()

	topic := domain_event.NotificationCreateTopic()

	go func() {
		if err := mq.SubscribeTo(mq.BackendRabbitMQ, topic, handleNotificationEvent); err != nil {
			logrus.Fatalf("Subscribe failed: %v", err)
		}
	}()

	utils.WaitForShutdownSignal()
	time.Sleep(500 * time.Millisecond)
	if err := mq.Close(); err != nil {
		logrus.Errorf("Close MQ failed: %v", err)
	}
	if err := storage.CloseRedis(); err != nil {
		logrus.Errorf("Close Redis failed: %v", err)
	}
	if err := storage.ClosePostgres(); err != nil {
		logrus.Errorf("Close Postgres failed: %v", err)
	}
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
