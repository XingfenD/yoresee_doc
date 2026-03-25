package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/dto"
	"github.com/XingfenD/yoresee_doc/internal/repository"
	"github.com/XingfenD/yoresee_doc/internal/service/notification_service"
	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

type notificationEvent struct {
	ReceiverExternalIDs []string        `json:"receiver_external_ids"`
	Type                string          `json:"type"`
	Title               string          `json:"title"`
	Content             string          `json:"content"`
	PayloadJSON         string          `json:"payload_json"`
	Payload             json.RawMessage `json:"payload"`
}

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

	topic := constant.NotificationTopicDefault

	mqBackend := strings.ToLower(os.Getenv("NOTIFICATION_MQ"))
	var backend mq.Backend
	switch mqBackend {
	case "rabbit", "rabbitmq":
		backend = mq.BackendRabbitMQ
	default:
		backend = mq.BackendRedis
	}

	logrus.Infof("Notification worker started: topic=%s backend=%s", topic, backend)

	go func() {
		if err := mq.SubscribeTo(backend, topic, handleNotificationEvent); err != nil {
			logrus.Fatalf("Subscribe failed: %v", err)
		}
	}()

	waitForShutdown()
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

	var evt notificationEvent
	if err := sonic.Unmarshal(data, &evt); err != nil {
		logrus.Errorf("parse notification event failed: %v", err)
		return nil
	}
	if evt.PayloadJSON == "" && len(evt.Payload) > 0 {
		evt.PayloadJSON = string(evt.Payload)
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

func waitForShutdown() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	time.Sleep(500 * time.Millisecond)
}
