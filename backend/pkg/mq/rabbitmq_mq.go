package mq

import (
	"context"
	"fmt"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

type RabbitMQConfig struct {
	URL string
}

func BuildRabbitMQConfig(cfg config.RabbitMQQueueConfig) RabbitMQConfig {
	return RabbitMQConfig{
		URL: cfg.URL,
	}
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

func (rmq *RabbitMQ) Publish(ctx context.Context, topic string, data []byte) error {
	err := rmq.channel.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	err = rmq.channel.PublishWithContext(
		ctx,
		topic,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

func (rmq *RabbitMQ) Subscribe(topic string, handler MessageHandler) error {
	err := rmq.channel.ExchangeDeclare(
		topic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange: %w", err)
	}

	q, err := rmq.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = rmq.channel.QueueBind(
		q.Name,
		"",
		topic,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind a queue: %w", err)
	}

	msgs, err := rmq.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			ctx := context.Background()
			if err := handler(ctx, d.Body); err != nil {
				fmt.Printf("Error handling message from topic %s: %v\n", topic, err)
			}
		}
	}()

	return nil
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
