package document_repo

import (
	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/redis/go-redis/v9"
)

type DocumentRepository struct {
	Loader cache_loader.Loader
}

var DocumentRepo DocumentRepository

func Init(redis *redis.Client) {
	DocumentRepo.Loader = *cache_loader.NewLoader(redis)
}
