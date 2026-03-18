package knowledge_base_repo

import (
	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/redis/go-redis/v9"
)

type KnowledgeBaseRepository struct {
	Loader cache_loader.Loader
}

var KnowledgeBaseRepo = &KnowledgeBaseRepository{}

func Init(redis *redis.Client) {
	KnowledgeBaseRepo.Loader = *cache_loader.NewLoader(redis)
}
