package mq

import (
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

func NewMessageQueue(mqType string, cfg *config.Config) (MessageQueue, error) {
	switch mqType {
	case "redis":
		return NewRedisMQ(), nil
	case "rabbitmq":
		rabbitCfg := RabbitMQConfig{
			URL: cfg.Backend.MQConfig.RabbitMQ.URL,
		}
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
