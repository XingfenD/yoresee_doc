package mq

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type Backend string

const (
	BackendRedis    Backend = "redis"
	BackendRabbitMQ Backend = "rabbitmq"
)

func NewMessageQueues(cfg *config.MQConfig) (map[Backend]MessageQueue, error) {
	mqs := map[Backend]MessageQueue{}

	mqs[BackendRedis] = NewRedisMQ()

	if cfg != nil && strings.TrimSpace(cfg.RabbitMQ.URL) != "" {
		rmq, err := NewRabbitMQ(RabbitMQConfig{URL: cfg.RabbitMQ.URL})
		if err != nil {
			return mqs, err
		}
		mqs[BackendRabbitMQ] = rmq
	}

	return mqs, nil
}
