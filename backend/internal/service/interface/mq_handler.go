package svc_iface

import "context"

type HandleFunc func(data []byte) error

type TopicConsumer interface {
	Topic() string
	Consume() HandleFunc
}

type TopicProducer[T any] interface {
	Publish(ctx context.Context, data T) error
	Handle(T) error
	Topic() string
}

// type MQTopicHandler[T any] interface {
// 	TopicProducer[T]
// 	TopicConsumer
// }
