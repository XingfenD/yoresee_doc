package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
)

type ConsulKVClient struct {
	baseURL    string
	token      string
	datacenter string
	prefix     string
	httpClient *http.Client
}

var Consul *ConsulKVClient

func InitConsul(cfg *config.ConsulConfig) error {
	if cfg == nil || !cfg.Enabled {
		Consul = nil
		return nil
	}

	scheme := strings.TrimSpace(cfg.Scheme)
	if scheme == "" {
		scheme = "http"
	}
	address := strings.TrimSpace(cfg.Address)
	if address == "" {
		address = "127.0.0.1:8500"
	}

	Consul = &ConsulKVClient{
		baseURL:    fmt.Sprintf("%s://%s", scheme, address),
		token:      cfg.Token,
		datacenter: cfg.Datacenter,
		prefix:     strings.Trim(cfg.Prefix, "/"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
	return nil
}

func ConsulEnabled() bool {
	return Consul != nil
}

func (c *ConsulKVClient) resolveKey(key string) string {
	key = strings.Trim(key, "/")
	if c.prefix == "" {
		return key
	}
	return path.Join(c.prefix, key)
}

func (c *ConsulKVClient) buildURL(key string, raw bool) string {
	fullKey := c.resolveKey(key)
	u, _ := url.Parse(c.baseURL)
	u.Path = path.Join(u.Path, "/v1/kv", fullKey)
	q := u.Query()
	if raw {
		q.Set("raw", "")
	}
	if c.datacenter != "" {
		q.Set("dc", c.datacenter)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *ConsulKVClient) setHeaders(req *http.Request) {
	if c.token != "" {
		req.Header.Set("X-Consul-Token", c.token)
	}
}

func (c *ConsulKVClient) Get(ctx context.Context, key string) (string, bool, error) {
	if c == nil {
		return "", false, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(key, true), nil)
	if err != nil {
		return "", false, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", false, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", false, fmt.Errorf("consul kv get failed: %s", string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, err
	}
	return string(body), true, nil
}

func (c *ConsulKVClient) Set(ctx context.Context, key, value string) error {
	if c == nil {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.buildURL(key, false), bytes.NewBufferString(value))
	if err != nil {
		return err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("consul kv set failed: %s", string(body))
	}
	return nil
}
