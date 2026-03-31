package mq

import "strings"

type Backend string

const (
	BackendRedis    Backend = "redis"
	BackendRabbitMQ Backend = "rabbitmq"
)

type Options struct {
	RabbitMQURL string
}

func NewMessageQueues(cfg *Options) (map[Backend]MessageQueue, error) {
	mqs := map[Backend]MessageQueue{}
	mqs[BackendRedis] = NewRedisMQ()

	if cfg != nil && strings.TrimSpace(cfg.RabbitMQURL) != "" {
		rmq, err := NewRabbitMQ(RabbitMQConfig{URL: cfg.RabbitMQURL})
		if err != nil {
			return mqs, err
		}
		mqs[BackendRabbitMQ] = rmq
	}

	return mqs, nil
}
