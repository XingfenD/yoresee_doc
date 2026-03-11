package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

func Init(cfg *config.MQConfig) error {
	mq, err := InitMessageQueue(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize message queue: %w", err)
	}

	MQ = mq

	return nil
}

func Close() error {
	return MQ.Close()
}

func Publish(ctx context.Context, topic string, data []byte) error {
	if MQ == nil {
		return fmt.Errorf("message queue not initialized")
	}
	return MQ.Publish(ctx, topic, data)
}

func Subscribe(topic string, handler MessageHandler) error {
	if MQ == nil {
		return fmt.Errorf("message queue not initialized")
	}
	return MQ.Subscribe(topic, handler)
}
