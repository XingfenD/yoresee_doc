package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type MessageHandler func(ctx context.Context, data []byte) error

type MessageQueue interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Close() error
}

var MQ MessageQueue

func NewMessageQueue(mqType string, cfg *config.Config) (MessageQueue, error) {
	switch mqType {
	case "redis":
		return NewRedisMQ(), nil
	case "rabbitmq":
		rabbitCfg := BuildRabbitMQConfig(cfg.Backend.MQConfig.RabbitMQ)
		return NewRabbitMQ(rabbitCfg)
	default:
		return nil, fmt.Errorf("unsupported mq type: %s", mqType)
	}
}

func InitMessageQueue(cfg *config.Config) (MessageQueue, error) {
	mqType := cfg.Backend.MQConfig.Type
	if mqType == "" {
		mqType = "redis"
	}

	return NewMessageQueue(mqType, cfg)
}
