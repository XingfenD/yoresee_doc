package service

import (
	"context"

	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

type MQService struct {
	mq mq.MessageQueue
}

var MQSvc = &MQService{
	mq: mq.MQ,
}

func (srvc *MQService) IsInitialized() bool {
	return srvc != nil && srvc.mq != nil
}

func (srvc *MQService) Publish(ctx context.Context, topic string, data []byte) error {
	if srvc == nil || srvc.mq == nil {
		return status.StatusMQNotInitialized
	}

	return MQSvc.mq.Publish(ctx, topic, data)
}

func (srvc *MQService) Subscribe(topic string, handler func([]byte) error) error {
	if srvc == nil || srvc.mq == nil {
		return status.StatusMQNotInitialized
	}

	return srvc.mq.Subscribe(topic, func(ctx context.Context, data []byte) error {
		return handler(data)
	})
}

func PublishByProducer[T any](ctx context.Context, producer svc_iface.TopicProducer, taskData T) error {
	if !MQSvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}
	data, err := sonic.Marshal(taskData)
	if err != nil {
		logrus.Error("marshal failed")
		return err
	}
	return MQSvc.Publish(ctx, producer.Topic(), data)
}
