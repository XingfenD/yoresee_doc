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
		ES = nil
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

func (c *ElasticsearchClient) SearchIDs(ctx context.Context, index string, searchBody map[string]interface{}) ([]int64, error) {
	if c == nil {
		return nil, fmt.Errorf("elasticsearch client is nil")
	}
	payload, err := json.Marshal(searchBody)
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

func (c *ElasticsearchClient) IndexExists(ctx context.Context, index string) (bool, error) {
	if c == nil {
		return false, fmt.Errorf("elasticsearch client is nil")
	}
	index = strings.TrimSpace(index)
	if index == "" {
		return false, fmt.Errorf("index is empty")
	}

	var lastErr error
	for _, address := range c.addresses {
		url := fmt.Sprintf("%s/%s", address, index)
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
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
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return true, nil
		}
		lastErr = fmt.Errorf("status=%d", resp.StatusCode)
	}
	return false, lastErr
}

func (c *ElasticsearchClient) CreateIndex(ctx context.Context, index string, body map[string]interface{}) error {
	if c == nil {
		return fmt.Errorf("elasticsearch client is nil")
	}
	index = strings.TrimSpace(index)
	if index == "" {
		return fmt.Errorf("index is empty")
	}

	payload := []byte("{}")
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = encoded
	}

	var lastErr error
	for _, address := range c.addresses {
		url := fmt.Sprintf("%s/%s", address, index)
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
		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		lastErr = fmt.Errorf("status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(bodyBytes)))
	}
	return lastErr
}

func CloseElasticsearch() error {
	ES = nil
	return nil
}
