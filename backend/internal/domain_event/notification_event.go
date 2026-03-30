package domain_event

import (
	"context"
	"strings"

	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/bytedance/sonic"
)

type NotificationCreateEvent struct {
	ReceiverExternalIDs []string `json:"receiver_external_ids"`
	Type                string   `json:"type"`
	Title               string   `json:"title"`
	Content             string   `json:"content"`
	PayloadJSON         string   `json:"payload_json"`
}

func NotificationCreateTopic() string {
	return constant.NotificationTopicDefault
}

func PublishNotificationCreateEvent(ctx context.Context, evt NotificationCreateEvent) error {
	return publishJSONToRabbitMQ(ctx, NotificationCreateTopic(), evt)
}

func DecodeNotificationCreateEvent(data []byte) (*NotificationCreateEvent, error) {
	var evt NotificationCreateEvent
	if err := sonic.Unmarshal(data, &evt); err != nil {
		return nil, err
	}
	cleanReceiverIDs := make([]string, 0, len(evt.ReceiverExternalIDs))
	for _, receiverID := range evt.ReceiverExternalIDs {
		receiverID = strings.TrimSpace(receiverID)
		if receiverID == "" {
			continue
		}
		cleanReceiverIDs = append(cleanReceiverIDs, receiverID)
	}
	evt.ReceiverExternalIDs = cleanReceiverIDs
	return &evt, nil
}
