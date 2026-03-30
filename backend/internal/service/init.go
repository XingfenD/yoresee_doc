package service

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service/config_service"
	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/internal/service/setting_service"
	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/sirupsen/logrus"
)

func RegisterTopicConsumer(h svc_iface.TopicConsumer) error {
	return mq_service.MQSvc.Consume(
		context.Background(),
		mq.BackendRedis,
		mq.ConsumeOptions{
			Topic:   h.Topic(),
			Mode:    mq.ConsumeModeFanout,
			AutoAck: true,
			OnError: mq.ErrorActionDrop,
		},
		func(ctx context.Context, message mq.Message) error {
			return h.Consume()(message.Body)
		},
	)

}

func Init(cfg *config.Config) error {
	config_service.InitConfigService()
	setting_service.InitSettingService()
	if err := InitMQTopicConsumer(); err != nil {
		logrus.Errorf("[Service layer] InitMQTopicConsumer failed, err=%+v", err)
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "initialize message queue topic consumer failed")
	}

	return nil
}

func InitMQTopicConsumer() error {
	return nil
}
