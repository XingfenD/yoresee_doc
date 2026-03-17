package cache

import (
	"context"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

func DoubleDelete(ctx context.Context, writeDB func() error, keys ...string) error {
	if err := deleteKeys(ctx, keys); err != nil {
		logrus.Warn("First cache deletion failed: ", err)
	}

	if err := writeDB(); err != nil {
		return err
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := deleteKeys(ctx, keys); err != nil {
			logrus.Warn("Second cache deletion failed: ", err)
		}
	}()

	return nil
}

func deleteKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	return storage.KVS.Del(ctx, keys...).Err()
}
