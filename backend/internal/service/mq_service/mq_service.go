package mq_service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/config"
	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/internal/utils"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/bytedance/sonic"
	"github.com/sirupsen/logrus"
)

type MQService struct {
	backend mq.Backend
	queues  map[mq.Backend]mq.MessageQueue
}

var MQSvc = &MQService{
	backend: mq.BackendRedis,
	queues:  map[mq.Backend]mq.MessageQueue{},
}

func (srvc *MQService) IsInitialized() bool {
	return srvc != nil && len(srvc.queues) > 0
}

func (srvc *MQService) Init(cfg *config.MQConfig) error {
	if srvc == nil {
		return status.StatusMQNotInitialized
	}

	queues, err := mq.NewMessageQueues(cfg)
	if err != nil {
		logrus.Errorf("[Service layer: MQService] Init failed, err=%+v", err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "init message queue failed")
	}
	srvc.queues = queues
	return nil
}

func (srvc *MQService) Close() error {
	if srvc == nil {
		return status.StatusMQNotInitialized
	}

	var firstErr error
	for backend, q := range srvc.queues {
		if q == nil {
			continue
		}
		if err := q.Close(); err != nil {
			logrus.Errorf("[Service layer: MQService] Close backend failed, backend=%s, err=%+v", backend, err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	srvc.queues = map[mq.Backend]mq.MessageQueue{}
	return firstErr
}

func (srvc *MQService) Publish(ctx context.Context, backend mq.Backend, msg mq.PublishMessage) error {
	if !srvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}

	msg.Topic = strings.TrimSpace(msg.Topic)
	if msg.Topic == "" {
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "publish message failed")
	}

	q, err := srvc.getQueue(backend)
	if err != nil {
		logrus.Errorf("[Service layer: MQService] Publish getQueue failed, backend=%s, err=%+v", backend, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "publish message failed")
	}

	if err := q.Publish(ctx, msg); err != nil {
		logrus.Errorf("[Service layer: MQService] Publish failed, backend=%s, topic=%s, err=%+v", backend, msg.Topic, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "publish message failed")
	}
	return nil
}

func (srvc *MQService) Consume(ctx context.Context, backend mq.Backend, opts mq.ConsumeOptions, handler func(context.Context, mq.Message) error) error {
	if !srvc.IsInitialized() {
		return status.StatusMQNotInitialized
	}
	opts = normalizeConsumeOptions(opts)
	if opts.Topic == "" {
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "subscribe message failed")
	}
	if err := validateConsumeOptions(opts); err != nil {
		logrus.Errorf("[Service layer: MQService] Consume options invalid, backend=%s, topic=%s, err=%+v", backend, opts.Topic, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "subscribe message failed")
	}

	q, err := srvc.getQueue(backend)
	if err != nil {
		logrus.Errorf("[Service layer: MQService] Consume getQueue failed, backend=%s, topic=%s, err=%+v", backend, opts.Topic, err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "subscribe message failed")
	}

	if err := q.Consume(ctx, opts, handler); err != nil {
		logrus.Errorf("[Service layer: MQService] Consume failed, backend=%s, topic=%s, err=%+v", backend, opts.Topic, err)
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
	return MQSvc.Publish(ctx, MQSvc.backend, mq.PublishMessage{
		Topic: producer.Topic(),
		Body:  data,
	})
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
	return MQSvc.Publish(ctx, backend, mq.PublishMessage{
		Topic: producer.Topic(),
		Body:  data,
	})
}

func normalizeConsumeOptions(opts mq.ConsumeOptions) mq.ConsumeOptions {
	n := opts
	n.Topic = strings.TrimSpace(n.Topic)
	n.Group = strings.TrimSpace(n.Group)
	n.Consumer = strings.TrimSpace(n.Consumer)
	if n.Mode == "" {
		n.Mode = mq.ConsumeModeFanout
	}
	if n.OnError == "" {
		n.OnError = mq.ErrorActionRequeue
	}

	if n.Mode == mq.ConsumeModeGroup {
		n.Group = utils.NormalizeToken(n.Group, "default")
	}

	consumer := n.Consumer
	if consumer == "" {
		hostName, _ := os.Hostname()
		consumer = fmt.Sprintf(
			"%s-%s-%d",
			utils.NormalizeToken(n.Topic, "default"),
			utils.NormalizeToken(hostName, "default"),
			os.Getpid(),
		)
	}
	n.Consumer = utils.NormalizeToken(consumer, "default")

	return n
}

func validateConsumeOptions(opts mq.ConsumeOptions) error {
	switch opts.Mode {
	case mq.ConsumeModeFanout, mq.ConsumeModeGroup:
	default:
		return fmt.Errorf("invalid consume mode: %s", opts.Mode)
	}

	switch opts.OnError {
	case mq.ErrorActionDrop, mq.ErrorActionRequeue:
	default:
		return fmt.Errorf("invalid error action: %s", opts.OnError)
	}

	return nil
}

func (srvc *MQService) getQueue(backend mq.Backend) (mq.MessageQueue, error) {
	if srvc == nil || len(srvc.queues) == 0 {
		return nil, fmt.Errorf("message queue not initialized")
	}
	q, ok := srvc.queues[backend]
	if !ok || q == nil {
		return nil, fmt.Errorf("message queue backend not initialized: %s", backend)
	}
	return q, nil
}
