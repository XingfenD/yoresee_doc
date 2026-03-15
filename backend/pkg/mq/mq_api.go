package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

func Init(cfg *config.MQConfig) error {
	mqs, err := InitMessageQueues(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize message queue: %w", err)
	}

	MQs = mqs

	return nil
}

func Close() error {
	for _, q := range MQs {
		if q != nil {
			_ = q.Close()
		}
	}
	return nil
}

func PublishTo(ctx context.Context, backend Backend, topic string, data []byte) error {
	if MQs == nil {
		return fmt.Errorf("message queue not initialized")
	}
	q, ok := MQs[backend]
	if !ok || q == nil {
		return fmt.Errorf("message queue backend not initialized: %s", backend)
	}
	return q.Publish(ctx, topic, data)
}

func SubscribeTo(backend Backend, topic string, handler MessageHandler) error {
	if MQs == nil {
		return fmt.Errorf("message queue not initialized")
	}
	q, ok := MQs[backend]
	if !ok || q == nil {
		return fmt.Errorf("message queue backend not initialized: %s", backend)
	}
	return q.Subscribe(topic, handler)
}
