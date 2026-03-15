package mq

import (
	"context"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type MessageHandler func(ctx context.Context, data []byte) error

type MessageQueue interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Close() error
}

type Backend string

const (
	BackendRedis    Backend = "redis"
	BackendRabbitMQ Backend = "rabbitmq"
)

var MQs = map[Backend]MessageQueue{}

func InitMessageQueues(cfg *config.MQConfig) (map[Backend]MessageQueue, error) {
	mqs := map[Backend]MessageQueue{}

	mqs[BackendRedis] = NewRedisMQ()

	if cfg != nil && cfg.RabbitMQ.URL != "" {
		rabbitCfg := BuildRabbitMQConfig(cfg.RabbitMQ)
		rmq, err := NewRabbitMQ(rabbitCfg)
		if err != nil {
			return mqs, err
		}
		mqs[BackendRabbitMQ] = rmq
	}

	return mqs, nil
}
