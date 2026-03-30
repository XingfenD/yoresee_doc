package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/redis/go-redis/v9"
)

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func StoreJWTToken(userExternalID, token string, ttl time.Duration) error {
	if storage.KVS == nil {
		return status.StatusRedisNotInitialized
	}
	if ttl <= 0 {
		return status.GenErrWithCustomMsg(status.StatusInternalParamsError, "invalid jwt ttl")
	}

	ctx := context.Background()
	tokenHash := hashToken(token)
	pipe := storage.KVS.TxPipeline()
	pipe.Set(ctx, cache.KeyJWTActiveToken(tokenHash), userExternalID, ttl)
	pipe.SAdd(ctx, cache.KeyJWTUserTokenSet(userExternalID), tokenHash)
	pipe.Expire(ctx, cache.KeyJWTUserTokenSet(userExternalID), ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("persist jwt token failed: %w", err)
	}
	return nil
}

func ValidateJWTTokenInRedis(userExternalID, token string) error {
	if storage.KVS == nil {
		return status.StatusRedisNotInitialized
	}

	ctx := context.Background()
	tokenHash := hashToken(token)
	blacklisted, err := storage.KVS.Exists(ctx, cache.KeyJWTBlacklistToken(tokenHash)).Result()
	if err != nil {
		return fmt.Errorf("query jwt blacklist failed: %w", err)
	}
	if blacklisted > 0 {
		return status.GenErrWithCustomMsg(status.StatusTokenInvalid, "token is blacklisted")
	}

	storedUserExternalID, err := storage.KVS.Get(ctx, cache.KeyJWTActiveToken(tokenHash)).Result()
	if err != nil {
		if err == redis.Nil {
			return status.GenErrWithCustomMsg(status.StatusTokenInvalid, "token is invalid")
		}
		return fmt.Errorf("query active jwt token failed: %w", err)
	}
	if storedUserExternalID != userExternalID {
		return status.GenErrWithCustomMsg(status.StatusTokenInvalid, "token user mismatch")
	}
	return nil
}

func BlacklistUserJWTs(userExternalID string) error {
	if storage.KVS == nil {
		return status.StatusRedisNotInitialized
	}

	ctx := context.Background()
	tokenHashes, err := storage.KVS.SMembers(ctx, cache.KeyJWTUserTokenSet(userExternalID)).Result()
	if err != nil {
		return fmt.Errorf("list user jwt tokens failed: %w", err)
	}
	if len(tokenHashes) == 0 {
		return nil
	}

	pipe := storage.KVS.TxPipeline()
	for _, tokenHash := range tokenHashes {
		activeKey := cache.KeyJWTActiveToken(tokenHash)
		ttlCmd := storage.KVS.TTL(ctx, activeKey)
		ttl, ttlErr := ttlCmd.Result()
		if ttlErr != nil && ttlErr != redis.Nil {
			return fmt.Errorf("query token ttl failed: %w", ttlErr)
		}
		if ttl <= 0 {
			ttl = 24 * time.Hour
		}
		pipe.Set(ctx, cache.KeyJWTBlacklistToken(tokenHash), "1", ttl)
		pipe.Del(ctx, activeKey)
	}
	pipe.Del(ctx, cache.KeyJWTUserTokenSet(userExternalID))
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("blacklist user jwt tokens failed: %w", err)
	}
	return nil
}
