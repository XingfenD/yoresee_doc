package domain_event

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/service/mq_service"
	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/bytedance/sonic"
)

func publishJSONToRabbitMQ(ctx context.Context, topic string, payload any) error {
	data, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}
	return mq_service.MQSvc.Publish(ctx, mq.BackendRabbitMQ, mq.PublishMessage{
		Topic: topic,
		Body:  data,
	})
}
