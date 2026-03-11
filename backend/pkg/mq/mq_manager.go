package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type MQManager struct {
	mq MessageQueue
}

var defaultManager *MQManager

func Init(cfg *config.Config) error {
	mq, err := InitMessageQueue(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize message queue: %w", err)
	}

	defaultManager = &MQManager{
		mq: mq,
	}

	return nil
}

func GetDefault() MessageQueue {
	if defaultManager == nil {
		return nil
	}
	return defaultManager.mq
}

func Close() error {
	if defaultManager != nil && defaultManager.mq != nil {
		return defaultManager.mq.Close()
	}
	return nil
}

func Publish(ctx context.Context, topic string, data []byte) error {
	mq := GetDefault()
	if mq == nil {
		return fmt.Errorf("message queue not initialized")
	}
	return mq.Publish(ctx, topic, data)
}

func Subscribe(topic string, handler MessageHandler) error {
	mq := GetDefault()
	if mq == nil {
		return fmt.Errorf("message queue not initialized")
	}
	return mq.Subscribe(topic, handler)
}
