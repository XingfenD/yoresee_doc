package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/config"
	"github.com/redis/go-redis/v9"
)

var KVS *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	KVS = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := KVS.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("init redis client failed: %w", err)
	}

	return nil
}

func GetRedis() *redis.Client {
	return KVS
}

func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return KVS.Set(ctx, key, value, expiration).Err()
}

func GetCache(ctx context.Context, key string) (string, error) {
	return KVS.Get(ctx, key).Result()
}

func DeleteCache(ctx context.Context, key string) error {
	return KVS.Del(ctx, key).Err()
}

func ClearCacheByPattern(ctx context.Context, pattern string) error {
	keys, err := KVS.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return KVS.Del(ctx, keys...).Err()
	}

	return nil
}
