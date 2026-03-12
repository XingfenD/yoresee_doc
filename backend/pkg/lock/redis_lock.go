package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/pkg/storage"
)

type DistributedLock struct {
	key        string
	value      string
	expiration time.Duration
}

func NewDistributedLock(key string, expiration time.Duration) *DistributedLock {
	return &DistributedLock{
		key:        key,
		value:      fmt.Sprintf("%d", time.Now().UnixNano()),
		expiration: expiration,
	}
}

func (dl *DistributedLock) Lock(ctx context.Context) error {
	success, err := storage.SetNX(ctx, dl.key, dl.value, dl.expiration)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to acquire lock: %s", dl.key)
	}
	return nil
}

func (dl *DistributedLock) TryLock(ctx context.Context) (bool, error) {
	return storage.SetNX(ctx, dl.key, dl.value, dl.expiration)
}

func (dl *DistributedLock) Unlock(ctx context.Context) error {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := storage.GetRedis().Eval(ctx, script, []string{dl.key}, dl.value).Result()
	if err != nil {
		return err
	}
	if result.(int64) == 0 {
		return fmt.Errorf("lock not held or expired: %s", dl.key)
	}
	return nil
}
