package domain_event

import (
	"context"
	"fmt"
	"strings"

	"github.com/XingfenD/yoresee_doc/pkg/constant"
	"github.com/bytedance/sonic"
)

const (
	DocumentActionUpsert = "upsert"
	DocumentActionDelete = "delete"
)

type DocumentSyncEvent struct {
	Action     string `json:"action"`
	ExternalID string `json:"external_id"`
}

func DocumentSyncTopic() string {
	return constant.SearchSyncTopicDefault
}

func PublishDocumentUpsertEvent(ctx context.Context, externalID string) error {
	externalID = strings.TrimSpace(externalID)
	if externalID == "" {
		return fmt.Errorf("external_id is empty")
	}
	return publishJSONToRabbitMQ(ctx, DocumentSyncTopic(), DocumentSyncEvent{
		Action:     DocumentActionUpsert,
		ExternalID: externalID,
	})
}

func DecodeDocumentSyncEvent(data []byte) (*DocumentSyncEvent, error) {
	var evt DocumentSyncEvent
	if err := sonic.Unmarshal(data, &evt); err != nil {
		return nil, err
	}
	evt.Action = strings.TrimSpace(evt.Action)
	evt.ExternalID = strings.TrimSpace(evt.ExternalID)
	if evt.Action == "" {
		evt.Action = DocumentActionUpsert
	}
	if evt.ExternalID == "" {
		return nil, fmt.Errorf("external_id is empty")
	}
	return &evt, nil
}
