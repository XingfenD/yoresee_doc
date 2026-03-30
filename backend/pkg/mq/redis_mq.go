package mq

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/errs"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisMQ struct {
	client *redis.Client
}

func NewRedisMQ() *RedisMQ {
	return &RedisMQ{
		client: storage.GetRedis(),
	}
}

func (rmq *RedisMQ) Publish(ctx context.Context, msg PublishMessage) error {
	topic := strings.TrimSpace(msg.Topic)
	if topic == "" {
		return errs.ErrTopicEmpty
	}
	if rmq.client == nil {
		return errs.ErrRedisClientNotInitialized
	}
	return rmq.client.Publish(ctx, topic, msg.Body).Err()
}

func (rmq *RedisMQ) Consume(ctx context.Context, opts ConsumeOptions, handler MessageHandler) error {
	opts = opts.normalize()
	topic := strings.TrimSpace(opts.Topic)
	if topic == "" {
		return errs.ErrTopicEmpty
	}
	if rmq.client == nil {
		return errs.ErrRedisClientNotInitialized
	}
	if opts.Mode == ConsumeModeGroup && strings.TrimSpace(opts.Group) != "" {
		logrus.Warnf("Redis Pub/Sub does not support true group consume, fallback to fanout mode, topic=%s, group=%s", topic, opts.Group)
	}

	pubsub := rmq.client.Subscribe(ctx, topic)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			if err := handler(context.Background(), Message{
				Topic:     topic,
				Body:      []byte(msg.Payload),
				Timestamp: time.Now(),
			}); err != nil {
				logrus.Errorf("Redis consume handler failed, topic=%s, err=%v", topic, err)
			}
		}
	}
}

func (rmq *RedisMQ) Close() error {
	return nil
}
