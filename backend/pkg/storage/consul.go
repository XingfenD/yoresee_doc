package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/XingfenD/yoresee_doc/pkg/errs"
)

type ConsulKVClient struct {
	baseURL    string
	token      string
	datacenter string
	prefix     string
	httpClient *http.Client
	cacheTTL   time.Duration
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
		cacheTTL:   5 * time.Second,
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

func (c *ConsulKVClient) CacheTTL() time.Duration {
	if c == nil || c.cacheTTL <= 0 {
		return 0
	}
	return c.cacheTTL
}

func (c *ConsulKVClient) ClearCache() {
	if c == nil {
		return
	}
	consulBindCache.Range(func(key, value any) bool {
		entry, ok := value.(*cacheEntry)
		if !ok {
			return true
		}
		entry.mu.Lock()
		entry.cache = cachedValue{}
		entry.mu.Unlock()
		return true
	})
}

func (c *ConsulKVClient) ClearCacheKey(key string) {
	if c == nil {
		return
	}
	val, ok := consulBindCache.Load(key)
	if !ok {
		return
	}
	entry, ok := val.(*cacheEntry)
	if !ok {
		return
	}
	entry.mu.Lock()
	entry.cache = cachedValue{}
	entry.mu.Unlock()
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
		return "", false, errs.Detail(errs.ErrConsulKVGet, string(body))
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
		return errs.Detail(errs.ErrConsulKVSet, string(body))
	}
	return nil
}

type consulTag struct {
	Key     string
	Default string
	IsJSON  bool
}

type cachedValue struct {
	value     string
	found     bool
	fetchedAt time.Time
}

type cacheEntry struct {
	mu    sync.Mutex
	cache cachedValue
}

var consulBindCache sync.Map

func BindConsulConfig(target interface{}, client *ConsulKVClient) error {
	if target == nil {
		return errs.ErrConsulBindTargetNil
	}
	if client == nil {
		return errs.ErrConsulClientNil
	}

	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errs.ErrConsulBindTargetInvalid
	}
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("consul")
		if tag == "" {
			continue
		}
		if field.Type.Kind() != reflect.Func || field.Type.NumIn() != 0 || field.Type.NumOut() != 1 {
			return errs.Detail(errs.ErrConsulFieldMustBeFunc, field.Name)
		}
		parsed, err := parseConsulTag(tag)
		if err != nil {
			return errs.DetailWrap(errs.ErrConsulFieldTagInvalid, field.Name, err)
		}

		outType := field.Type.Out(0)
		cache := &cachedValue{}
		entry := &cacheEntry{cache: *cache}
		consulBindCache.Store(parsed.Key, entry)

		fn := reflect.MakeFunc(field.Type, func(args []reflect.Value) []reflect.Value {
			entry.mu.Lock()
			defer entry.mu.Unlock()

			now := time.Now()
			ttl := client.CacheTTL()
			if ttl > 0 && entry.cache.fetchedAt.Add(ttl).After(now) {
				val, err := convertConsulValue(entry.cache.value, entry.cache.found, parsed, outType)
				if err == nil {
					return []reflect.Value{val}
				}
			}

			raw, found, err := client.Get(context.Background(), parsed.Key)
			if err != nil {
				val, _ := convertConsulValue(entry.cache.value, entry.cache.found, parsed, outType)
				return []reflect.Value{val}
			}
			entry.cache.value = raw
			entry.cache.found = found
			entry.cache.fetchedAt = now

			val, err := convertConsulValue(raw, found, parsed, outType)
			if err != nil {
				val, _ = convertConsulValue(entry.cache.value, entry.cache.found, parsed, outType)
			}
			return []reflect.Value{val}
		})

		val.Field(i).Set(fn)
	}

	return nil
}

func parseConsulTag(tag string) (consulTag, error) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return consulTag{}, errs.ErrConsulTagEmpty
	}
	parts := strings.Split(tag, ",")
	result := consulTag{
		Key: strings.TrimSpace(parts[0]),
	}
	if result.Key == "" {
		return consulTag{}, errs.ErrConsulTagMissingKey
	}
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if part == "json" {
			result.IsJSON = true
			continue
		}
		if strings.HasPrefix(part, "default=") {
			result.Default = strings.TrimPrefix(part, "default=")
		}
	}
	return result, nil
}

func convertConsulValue(raw string, found bool, tag consulTag, outType reflect.Type) (reflect.Value, error) {
	if !found || strings.TrimSpace(raw) == "" {
		raw = tag.Default
	}
	if tag.IsJSON || outType.Kind() == reflect.Struct || outType.Kind() == reflect.Map || outType.Kind() == reflect.Slice {
		target := reflect.New(outType).Interface()
		if err := json.Unmarshal([]byte(raw), target); err != nil {
			return reflect.Zero(outType), err
		}
		return reflect.ValueOf(target).Elem(), nil
	}

	switch outType.Kind() {
	case reflect.String:
		return reflect.ValueOf(raw).Convert(outType), nil
	case reflect.Bool:
		val, err := strconv.ParseBool(strings.TrimSpace(raw))
		if err != nil {
			return reflect.Zero(outType), err
		}
		return reflect.ValueOf(val).Convert(outType), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return reflect.Zero(outType), err
		}
		return reflect.ValueOf(val).Convert(outType), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return reflect.Zero(outType), err
		}
		return reflect.ValueOf(val).Convert(outType), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
		if err != nil {
			return reflect.Zero(outType), err
		}
		return reflect.ValueOf(val).Convert(outType), nil
	default:
		return reflect.Zero(outType), errs.Detailf(errs.ErrConsulUnsupportedType, "%s", outType.Kind())
	}
}
