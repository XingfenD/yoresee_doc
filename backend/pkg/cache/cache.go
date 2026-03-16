package cache

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

func GetJSON(ctx context.Context, key string, dst interface{}) (bool, error) {
	b, err := storage.KVS.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := sonic.Unmarshal(b, dst); err != nil {
		return false, err
	}
	return true, nil
}

func SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	b, err := sonic.Marshal(v)
	if err != nil {
		return err
	}
	return storage.KVS.Set(ctx, key, b, ttl).Err()
}

func DeleteByPattern(ctx context.Context, pattern string) error {
	iter := storage.KVS.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		_ = storage.KVS.Del(ctx, iter.Val()).Err()
	}
	return iter.Err()
}
