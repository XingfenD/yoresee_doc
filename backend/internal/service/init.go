package service

import (
	"github.com/XingfenD/yoresee_doc/internal/config"
	svc_iface "github.com/XingfenD/yoresee_doc/internal/service/interface"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
)

func RegisterTopicConsumer(h svc_iface.TopicConsumer) error {
	return MQSvc.SubscribeTo(mq.BackendRedis, h.Topic(), h.Consume())

}

func Init(cfg *config.Config) error {
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
