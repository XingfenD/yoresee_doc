package service

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/service/config_service"
	"github.com/XingfenD/yoresee_doc/internal/service/setting_service"
	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
)

func RegisterTopicConsumer(h svc_iface.TopicConsumer) error {
	return mq_service.MQSvc.SubscribeTo(mq.BackendRedis, h.Topic(), h.Consume())

}

func Init(cfg *config.Config) error {
	config_service.InitConfigService()
	setting_service.InitSettingService()
	// by cluster init the consumer (only init on consumer cluster)
	clusterRole := cfg.Backend.ClusterRole
	if clusterRole == config.ClusterRoleConsumer ||
		clusterRole == config.ClusterRoleHybrid {
		if err := InitMQTopicConsumer(); err != nil {
			return err
		}
	}

	return nil
}

func InitMQTopicConsumer() error {
	return nil
}
