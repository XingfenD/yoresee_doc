package mq

import "context"

type MessageQueue interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Close() error
}

type MessageHandler func(ctx context.Context, data []byte) error
