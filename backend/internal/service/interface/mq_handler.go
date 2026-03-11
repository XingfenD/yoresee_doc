package svc_iface

type HandleFunc func(data []byte) error

type TopicConsumer interface {
	Topic() string
	Consume() HandleFunc
}

type TopicProducer interface {
	Topic() string
}

// type MQTopicHandler interface {
// 	TopicProducer
// 	TopicConsumer
// }
