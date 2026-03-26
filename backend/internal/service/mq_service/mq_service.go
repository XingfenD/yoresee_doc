package mq_service

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

	if err := mq.PublishTo(ctx, backend, topic, data); err != nil {
		logrus.Errorf("[Service layer: MQService] PublishTo failed, backend=%s, topic=%s, err=%+v", backend, topic, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "publish message failed")
	}
	return nil
}

func (srvc *MQService) SubscribeTo(backend mq.Backend, topic string, handler func([]byte) error) error {
	if srvc == nil {
		return status.StatusMQNotInitialized
	}

	if err := mq.SubscribeTo(backend, topic, func(ctx context.Context, data []byte) error {
		return handler(data)
	}); err != nil {
		logrus.Errorf("[Service layer: MQService] SubscribeTo failed, backend=%s, topic=%s, err=%+v", backend, topic, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "subscribe message failed")
	}
	return nil
}

func PublishByProducer[T any](ctx context.Context, producer svc_iface.TopicProducer, taskData T) error {
	if !MQSvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}
	data, err := sonic.Marshal(taskData)
	if err != nil {
		logrus.Errorf("[Service layer: MQService] marshal failed, producer=%s, err=%+v", producer.Topic(), err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "marshal message failed")
	}
	return MQSvc.PublishTo(ctx, MQSvc.backend, producer.Topic(), data)
}

func PublishByProducerTo[T any](ctx context.Context, backend mq.Backend, producer svc_iface.TopicProducer, taskData T) error {
	if !MQSvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}
	data, err := sonic.Marshal(taskData)
	if err != nil {
		logrus.Errorf("[Service layer: MQService] marshal failed, backend=%s, producer=%s, err=%+v", backend, producer.Topic(), err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "marshal message failed")
	}
	return MQSvc.PublishTo(ctx, backend, producer.Topic(), data)
}
