package lock

import (
	"context"
	"time"
)

const LockCachePrefix = "lock:"

type RetryLockFunc func(ctx context.Context) (interface{}, error)

type CheckFunc func(ctx context.Context) (result interface{}, completed bool, err error)

func AcquireWithRetry(ctx context.Context, lockKey string, expiration time.Duration,
	maxRetries int, retryInterval time.Duration, checkFn CheckFunc, execFn RetryLockFunc) (interface{}, error) {

	distributedLock := NewDistributedLock(lockKey, expiration)

	err := distributedLock.Lock(ctx)
	if err == nil {
		defer distributedLock.Unlock(ctx)

		return execFn(ctx)
	}

	for i := 0; i < maxRetries; i++ {
		time.Sleep(retryInterval)

		if checkFn != nil {
			result, completed, err := checkFn(ctx)
			if err == nil && completed {
				return result, nil
			}
		}

		acquired, lockErr := distributedLock.TryLock(ctx)
		if lockErr != nil {
			continue
		}

		if acquired {
			defer distributedLock.Unlock(ctx)

			return execFn(ctx)
		}
	}

	return nil, err
}
