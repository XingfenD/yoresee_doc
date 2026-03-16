package cache

import (
	"encoding/json"

	"github.com/XingfenD/yoresee_doc/internal/model"
)

type Parser[T any] func(data []byte) (*T, error)

func DefaultParser[T any](data []byte) (*T, error) {
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

var ParseInt64 = DefaultParser[int64]

func ParseIDFromDocument(data []byte) (*int64, error) {
	var doc model.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc.ID, nil
}
