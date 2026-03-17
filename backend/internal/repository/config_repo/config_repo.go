package config_repo

import (
	"time"

	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/redis/go-redis/v9"
)

type ConfigRepository struct {
	Loader cache_loader.Loader
}

var ConfigRepo = ConfigRepository{}

const (
	cacheExpiration = 7 * 24 * time.Hour
)

func Init(redis *redis.Client) {
	ConfigRepo.Loader = *cache_loader.NewLoader(redis)
}
