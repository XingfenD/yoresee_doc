package cache

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type Loader struct {
	redis *redis.Client
	sf    singleflight.Group
}

func (l *Loader) Redis() *redis.Client {
	return l.redis
}

func NewLoader(redis *redis.Client) *Loader {
	return &Loader{
		redis: redis,
	}
}

type CacheLoadOperation[T any] struct {
	loader       *Loader
	keyParserMap map[string]Parser[T]
	defaultKey   *string
	baseTTL      time.Duration
	dbLoader     func() (*T, error)
	valueType    bool
}

func NewCacheLoadOperation[T any](loader *Loader) *CacheLoadOperation[T] {
	return &CacheLoadOperation[T]{
		loader:       loader,
		baseTTL:      time.Hour,
		valueType:    false,
		keyParserMap: make(map[string]Parser[T]),
	}
}

func (op *CacheLoadOperation[T]) WithDefaultKeyAndParser(key string, parser Parser[T]) *CacheLoadOperation[T] {
	var defaultParser Parser[T]
	op.defaultKey = &key
	if parser == nil {
		defaultParser = DefaultParser[T]
	} else {
		defaultParser = parser
	}

	op.keyParserMap[key] = defaultParser
	return op
}

func (op *CacheLoadOperation[T]) WithKeyAndParser(key string, parser Parser[T]) *CacheLoadOperation[T] {
	op.keyParserMap[key] = parser
	return op
}

func (op *CacheLoadOperation[T]) WithTTL(ttl time.Duration) *CacheLoadOperation[T] {
	op.baseTTL = ttl
	return op
}

func (op *CacheLoadOperation[T]) WithDBLoader(dbLoader func() (*T, error)) *CacheLoadOperation[T] {
	op.dbLoader = dbLoader
	return op
}

func (op *CacheLoadOperation[T]) Exec(ctx context.Context) (*T, error) {
	if op.defaultKey == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusInternalParamsError, "default key is needed")
	}
	if op.dbLoader == nil {
		return nil, status.GenErrWithCustomMsg(status.StatusInternalParamsError, "database loader function is required")
	}

	// try to load from any of the cache keys
	for key, parser := range op.keyParserMap {
		val, err := op.loader.redis.Get(ctx, key).Result()
		if err == nil {
			if val == "null" {
				return nil, nil
			}

			if data, err := parser([]byte(val)); err == nil {
				return data, nil
			}
		}
	}

	// all caches missed, use singleflight to avoid concurrent db queries
	sfKey := *op.defaultKey
	v, err, _ := op.loader.sf.Do(sfKey, func() (interface{}, error) {
		// double check all caches
		for key, parser := range op.keyParserMap {

			val, err := op.loader.redis.Get(ctx, key).Result()
			if err == nil {
				if val == "null" {
					return nil, nil
				}

				if data, err := parser([]byte(val)); err == nil {
					return data, nil
				}
			}
		}

		// query database
		data, err := op.dbLoader()
		if err != nil {
			// cache empty value for all keys if record not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				op.loader.redis.Set(ctx, *op.defaultKey, "null", 5*time.Minute)
			}
			return nil, err
		}

		bytes, _ := sonic.Marshal(data)
		ttl := op.baseTTL + time.Duration(rand.Intn(600))*time.Second

		op.loader.redis.Set(ctx, *op.defaultKey, bytes, ttl)

		return data, nil
	})

	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	return v.(*T), nil
}
