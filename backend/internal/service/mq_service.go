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
	backend mq.Backend
}

var MQSvc = &MQService{
	backend: mq.BackendRedis,
}

func (srvc *MQService) IsInitialized() bool {
	return srvc != nil
}

func (srvc *MQService) PublishTo(ctx context.Context, backend mq.Backend, topic string, data []byte) error {
	if srvc == nil {
		return status.StatusMQNotInitialized
	}

	return mq.PublishTo(ctx, backend, topic, data)
}

func (srvc *MQService) SubscribeTo(backend mq.Backend, topic string, handler func([]byte) error) error {
	if srvc == nil {
		return status.StatusMQNotInitialized
	}

	return mq.SubscribeTo(backend, topic, func(ctx context.Context, data []byte) error {
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
	return MQSvc.PublishTo(ctx, MQSvc.backend, producer.Topic(), data)
}

func PublishByProducerTo[T any](ctx context.Context, backend mq.Backend, producer svc_iface.TopicProducer, taskData T) error {
	if !MQSvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}
	data, err := sonic.Marshal(taskData)
	if err != nil {
		logrus.Error("marshal failed")
		return err
	}
	return MQSvc.PublishTo(ctx, backend, producer.Topic(), data)
}
