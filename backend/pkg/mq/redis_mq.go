package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/redis/go-redis/v9"
)

type RedisMQ struct {
	client *redis.Client
}

func NewRedisMQ() *RedisMQ {
	return &RedisMQ{
		client: storage.GetRedis(),
	}
}

func (rmq *RedisMQ) Publish(ctx context.Context, topic string, data []byte) error {
	return rmq.client.Publish(ctx, topic, data).Err()
}

func (rmq *RedisMQ) Subscribe(topic string, handler MessageHandler) error {
	pubsub := rmq.client.Subscribe(context.Background(), topic)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		ctx := context.Background()
		if err := handler(ctx, []byte(msg.Payload)); err != nil {
			fmt.Printf("Error handling message from topic %s: %v\n", topic, err)
		}
	}
	return nil
}

func (rmq *RedisMQ) Close() error {
	return nil
}
