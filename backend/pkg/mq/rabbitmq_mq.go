package mq

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

type RabbitMQConfig struct {
	URL string
}

func NewRabbitMQ(config RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (rmq *RabbitMQ) Publish(ctx context.Context, msg PublishMessage) error {
	topic := strings.TrimSpace(msg.Topic)
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}
	if err := rmq.declareTopicExchange(topic); err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	if err := rmq.channel.PublishWithContext(
		ctx,
		topic,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msg.Body,
			MessageId:   strings.TrimSpace(msg.Key),
			Headers:     toAMQPHeaders(msg.Headers),
			Timestamp:   time.Now(),
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

func (rmq *RabbitMQ) Consume(ctx context.Context, opts ConsumeOptions, handler MessageHandler) error {
	opts = opts.normalize()
	topic := strings.TrimSpace(opts.Topic)
	if topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if err := rmq.declareTopicExchange(topic); err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	queueName, durable, autoDelete, exclusive := consumeQueueConfig(topic, opts)
	queue, err := rmq.channel.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	if err := rmq.channel.QueueBind(
		queue.Name,
		"",
		topic,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	consumerTag := buildConsumerTag(opts.Consumer, topic)
	msgs, err := rmq.channel.Consume(
		queue.Name,
		consumerTag,
		opts.AutoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	logrus.Infof("RabbitMQ consume started, topic=%s, queue=%s, mode=%s, group=%s, auto_ack=%v",
		topic, queue.Name, opts.Mode, opts.Group, opts.AutoAck,
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("consumer channel closed, topic=%s, queue=%s", topic, queue.Name)
			}

			message := Message{
				ID:        d.MessageId,
				Topic:     topic,
				Key:       d.MessageId,
				Body:      d.Body,
				Headers:   fromAMQPHeaders(d.Headers),
				Timestamp: d.Timestamp,
			}
			if message.Timestamp.IsZero() {
				message.Timestamp = time.Now()
			}

			handleErr := handler(context.Background(), message)
			if opts.AutoAck {
				if handleErr != nil {
					logrus.Errorf("RabbitMQ consume handler failed (auto-ack), topic=%s, queue=%s, err=%v", topic, queue.Name, handleErr)
				}
				continue
			}

			if handleErr != nil {
				requeue := opts.OnError == ErrorActionRequeue
				if nackErr := d.Nack(false, requeue); nackErr != nil {
					logrus.Errorf("RabbitMQ nack failed, topic=%s, queue=%s, requeue=%v, err=%v", topic, queue.Name, requeue, nackErr)
				}
				logrus.Errorf("RabbitMQ consume handler failed, topic=%s, queue=%s, requeue=%v, err=%v", topic, queue.Name, requeue, handleErr)
				continue
			}

			if ackErr := d.Ack(false); ackErr != nil {
				logrus.Errorf("RabbitMQ ack failed, topic=%s, queue=%s, err=%v", topic, queue.Name, ackErr)
			}
		}
	}
}

func (rmq *RabbitMQ) declareTopicExchange(topic string) error {
	return rmq.channel.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func consumeQueueConfig(topic string, opts ConsumeOptions) (queueName string, durable, autoDelete, exclusive bool) {
	if opts.Mode != ConsumeModeGroup {
		return "", false, true, true
	}
	group := strings.TrimSpace(opts.Group)
	if group == "" {
		group = "default"
	}
	return fmt.Sprintf("%s.%s", topic, group), true, false, false
}

func buildConsumerTag(consumer, topic string) string {
	consumer = strings.TrimSpace(consumer)
	if consumer == "" {
		hostName, _ := os.Hostname()
		consumer = fmt.Sprintf("%s-%s-%d", strings.TrimSpace(topic), strings.TrimSpace(hostName), os.Getpid())
	}
	return consumer
}

func toAMQPHeaders(headers map[string]string) amqp091.Table {
	if len(headers) == 0 {
		return nil
	}
	table := amqp091.Table{}
	for key, val := range headers {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		table[key] = val
	}
	if len(table) == 0 {
		return nil
	}
	return table
}

func fromAMQPHeaders(table amqp091.Table) map[string]string {
	if len(table) == 0 {
		return nil
	}
	headers := make(map[string]string, len(table))
	for key, val := range table {
		headers[key] = fmt.Sprint(val)
	}
	return headers
}

func (rmq *RabbitMQ) Close() error {
	if rmq.channel != nil {
		rmq.channel.Close()
	}
	if rmq.conn != nil {
		rmq.conn.Close()
	}
	return nil
}
