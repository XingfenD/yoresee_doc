package domain_event

import (
	"context"

	"github.com/XingfenD/yoresee_doc/pkg/mq"
	"github.com/bytedance/sonic"
)

func publishJSONToRabbitMQ(ctx context.Context, topic string, payload any) error {
	data, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}
	return mq.PublishTo(ctx, mq.BackendRabbitMQ, topic, data)
}
