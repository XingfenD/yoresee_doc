package main

import (
	"strings"

	"github.com/bytedance/sonic"
)

type dirtyDocMessage struct {
	DocID string `json:"doc_id"`
	DocId string `json:"docId"`
}

func parseDocID(data []byte) string {
	payload := strings.TrimSpace(string(data))
	if payload == "" {
		return ""
	}
	var msg dirtyDocMessage
	if err := sonic.Unmarshal(data, &msg); err == nil {
		if msg.DocID != "" {
			return msg.DocID
		}
		if msg.DocId != "" {
			return msg.DocId
		}
	}
	return payload
}
