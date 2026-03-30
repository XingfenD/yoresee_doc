package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type ElasticsearchClient struct {
	addresses []string
	username  string
	password  string
	client    *http.Client
}

var ES *ElasticsearchClient

type SearchDocumentsRequest struct {
	Keyword              string
	UserID               *int64
	KnowledgeID          *int64
	ListOwnDoc           bool
	DocType              *string
	Status               *int
	Tags                 []string
	CreateTimeRangeStart *string
	CreateTimeRangeEnd   *string
	UpdateTimeRangeStart *string
	UpdateTimeRangeEnd   *string
	Size                 int
}

func InitElasticsearch(cfg *config.ElasticsearchConfig) error {
	if cfg == nil || !cfg.Enabled {
		ES = nil
		return nil
	}
	addresses := make([]string, 0, len(cfg.Addresses))
	for _, addr := range cfg.Addresses {
		trimmed := strings.TrimSpace(addr)
		if trimmed == "" {
			continue
		}
		addresses = append(addresses, strings.TrimRight(trimmed, "/"))
	}
	if len(addresses) == 0 {
		return fmt.Errorf("elasticsearch addresses are empty")
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	ES = &ElasticsearchClient{
		addresses: addresses,
		username:  cfg.Username,
		password:  cfg.Password,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	if err := ES.Ping(context.Background()); err != nil {
		return fmt.Errorf("ping elasticsearch failed: %w", err)
	}
	return nil
}

func (c *ElasticsearchClient) Ping(ctx context.Context) error {
	if c == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
	var lastErr error
	for _, address := range c.addresses {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
		if err != nil {
			lastErr = err
			continue
		}
		if c.username != "" {
			req.SetBasicAuth(c.username, c.password)
		}
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body := struct {
			Tagline string `json:"tagline"`
		}{}
		decodeErr := json.NewDecoder(resp.Body).Decode(&body)
		_ = resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 && decodeErr == nil {
			return nil
		}
		lastErr = fmt.Errorf("status=%d decodeErr=%v", resp.StatusCode, decodeErr)
	}
	return lastErr
}

func (c *ElasticsearchClient) UpsertDocument(ctx context.Context, index string, docID string, body map[string]interface{}) error {
	if c == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
	if strings.TrimSpace(index) == "" || strings.TrimSpace(docID) == "" {
		return fmt.Errorf("index or docID is empty")
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	var lastErr error
	for _, address := range c.addresses {
		url := fmt.Sprintf("%s/%s/_doc/%s", address, index, docID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(payload))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		if c.username != "" {
			req.SetBasicAuth(c.username, c.password)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		lastErr = fmt.Errorf("status=%d", resp.StatusCode)
	}
	return lastErr
}

func (c *ElasticsearchClient) SearchDocumentIDs(ctx context.Context, index string, req SearchDocumentsRequest) ([]int64, error) {
	if c == nil {
		return nil, fmt.Errorf("elasticsearch client is nil")
	}
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
	if req.Status != nil {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{"status": *req.Status},
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

	searchReq := map[string]interface{}{
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

	payload, err := json.Marshal(searchReq)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for _, address := range c.addresses {
		url := fmt.Sprintf("%s/%s/_search", address, index)
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			lastErr = err
			continue
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if c.username != "" {
			httpReq.SetBasicAuth(c.username, c.password)
		}

		resp, err := c.client.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
			continue
		}

		var parsed struct {
			Hits struct {
				Hits []struct {
					ID     string `json:"_id"`
					Source struct {
						ID int64 `json:"id"`
					} `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		}
		if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
			lastErr = err
			continue
		}

		ids := make([]int64, 0, len(parsed.Hits.Hits))
		for _, hit := range parsed.Hits.Hits {
			if hit.Source.ID > 0 {
				ids = append(ids, hit.Source.ID)
				continue
			}
			id, err := strconv.ParseInt(hit.ID, 10, 64)
			if err != nil || id <= 0 {
				continue
			}
			ids = append(ids, id)
		}
		return ids, nil
	}
	return nil, lastErr
}

func CloseElasticsearch() error {
	ES = nil
	return nil
}
