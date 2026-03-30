package mq

import (
	"context"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/errs"
)

type Message struct {
	ID        string
	Topic     string
	Key       string
	Body      []byte
	Headers   map[string]string
	Timestamp time.Time
}

type PublishMessage struct {
	Topic   string
	Key     string
	Body    []byte
	Headers map[string]string
}

type ConsumeMode string

const (
	ConsumeModeFanout ConsumeMode = "fanout"
	ConsumeModeGroup  ConsumeMode = "group"
)

type ErrorAction string

const (
	ErrorActionDrop    ErrorAction = "drop"
	ErrorActionRequeue ErrorAction = "requeue"
)

type ConsumeOptions struct {
	Topic    string
	Mode     ConsumeMode
	Group    string
	Consumer string
	AutoAck  bool
	OnError  ErrorAction
}

func (o ConsumeOptions) normalize() ConsumeOptions {
	n := o
	n.Topic = strings.TrimSpace(n.Topic)
	n.Group = strings.TrimSpace(n.Group)
	n.Consumer = strings.TrimSpace(n.Consumer)
	if n.Mode == "" {
		n.Mode = ConsumeModeFanout
	}
	if n.OnError == "" {
		n.OnError = ErrorActionRequeue
	}
	return n
}

func (o ConsumeOptions) validate() error {
	switch o.Mode {
	case ConsumeModeFanout, ConsumeModeGroup:
	default:
		return errs.Detailf(errs.ErrInvalidConsumeMode, "%s", o.Mode)
	}

	switch o.OnError {
	case ErrorActionDrop, ErrorActionRequeue:
	default:
		return errs.Detailf(errs.ErrInvalidErrorAction, "%s", o.OnError)
	}

	return nil
}

type MessageHandler func(ctx context.Context, msg Message) error

type MessageQueue interface {
	Publish(ctx context.Context, msg PublishMessage) error
	Consume(ctx context.Context, opts ConsumeOptions, handler MessageHandler) error
	Close() error
}
