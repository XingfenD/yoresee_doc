package search

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

const defaultSearchIndexPrefix = "yoresee_doc"

type DocumentSearchRequest struct {
	Keyword              string
	UserID               *int64
	KnowledgeID          *int64
	ListOwnDoc           bool
	DocType              *string
	Tags                 []string
	CreateTimeRangeStart *string
	CreateTimeRangeEnd   *string
	UpdateTimeRangeStart *string
	UpdateTimeRangeEnd   *string
	Size                 int
}

func SearchDocumentIDs(ctx context.Context, req DocumentSearchRequest) ([]int64, error) {
	if storage.ES == nil {
		return nil, fmt.Errorf("elasticsearch client is nil")
	}
	return storage.ES.SearchIDs(ctx, DocumentIndexName(), BuildDocumentSearchBody(req))
}

func UpsertDocument(ctx context.Context, doc *model.Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	if storage.ES == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
	return storage.ES.UpsertDocument(ctx, DocumentIndexName(), strconv.FormatInt(doc.ID, 10), BuildDocumentIndexBody(doc))
}

func DocumentIndexName() string {
	prefix := defaultSearchIndexPrefix
	if config.GlobalConfig != nil {
		configPrefix := strings.TrimSpace(config.GlobalConfig.Elasticsearch.IndexPrefix)
		if configPrefix != "" {
			prefix = configPrefix
		}
	}
	return key.KeyESDocumentIndex(prefix)
}

func BuildDocumentSearchBody(req DocumentSearchRequest) map[string]interface{} {
	size := req.Size
	if size <= 0 {
		size = 5000
	}
	if size > 10000 {
		size = 10000
	}

	must := []map[string]interface{}{}
	if strings.TrimSpace(req.Keyword) != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     strings.TrimSpace(req.Keyword),
				"fields":    []string{"title^3", "summary^2", "content"},
				"type":      "best_fields",
				"operator":  "and",
				"fuzziness": "AUTO",
			},
		})
	}

	filter := []map[string]interface{}{}
	mustNot := []map[string]interface{}{}
	if req.UserID != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"user_id": *req.UserID},
		})
	}
	if req.ListOwnDoc {
		mustNot = append(mustNot, map[string]interface{}{
			"exists": map[string]interface{}{"field": "knowledge_id"},
		})
	}
	if req.KnowledgeID != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"knowledge_id": *req.KnowledgeID},
		})
	}
	if req.DocType != nil && strings.TrimSpace(*req.DocType) != "" {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"type": strings.TrimSpace(*req.DocType)},
		})
	}
	for _, tag := range req.Tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"tags.keyword": tag},
		})
	}

	rangeCreated := map[string]interface{}{}
	if req.CreateTimeRangeStart != nil && strings.TrimSpace(*req.CreateTimeRangeStart) != "" {
		rangeCreated["gte"] = strings.TrimSpace(*req.CreateTimeRangeStart)
	}
	if req.CreateTimeRangeEnd != nil && strings.TrimSpace(*req.CreateTimeRangeEnd) != "" {
		rangeCreated["lte"] = strings.TrimSpace(*req.CreateTimeRangeEnd)
	}
	if len(rangeCreated) > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{"created_at": rangeCreated},
		})
	}

	rangeUpdated := map[string]interface{}{}
	if req.UpdateTimeRangeStart != nil && strings.TrimSpace(*req.UpdateTimeRangeStart) != "" {
		rangeUpdated["gte"] = strings.TrimSpace(*req.UpdateTimeRangeStart)
	}
	if req.UpdateTimeRangeEnd != nil && strings.TrimSpace(*req.UpdateTimeRangeEnd) != "" {
		rangeUpdated["lte"] = strings.TrimSpace(*req.UpdateTimeRangeEnd)
	}
	if len(rangeUpdated) > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{"updated_at": rangeUpdated},
		})
	}

	return map[string]interface{}{
		"size": size,
		"_source": []string{
			"id",
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":     must,
				"filter":   filter,
				"must_not": mustNot,
			},
		},
	}
}

func BuildDocumentIndexBody(doc *model.Document) map[string]interface{} {
	payload := map[string]interface{}{
		"id":          doc.ID,
		"external_id": doc.ExternalID,
		"title":       doc.Title,
		"summary":     doc.Summary,
		"content":     doc.Content,
		"type":        doc.Type,
		"user_id":     doc.UserID,
		"tags":        doc.Tags,
		"created_at":  doc.CreatedAt.Format(time.RFC3339),
		"updated_at":  doc.UpdatedAt.Format(time.RFC3339),
	}
	if doc.KnowledgeID != nil {
		payload["knowledge_id"] = *doc.KnowledgeID
	}
	return payload
}

func BuildDocumentIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type": "long",
				},
				"external_id": map[string]interface{}{
					"type": "keyword",
				},
				"title": map[string]interface{}{
					"type": "text",
				},
				"summary": map[string]interface{}{
					"type": "text",
				},
				"content": map[string]interface{}{
					"type": "text",
				},
				"type": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "long",
				},
				"knowledge_id": map[string]interface{}{
					"type": "long",
				},
				"tags": map[string]interface{}{
					"type": "keyword",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
			},
		},
	}
}
