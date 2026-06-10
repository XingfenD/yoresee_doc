package knowledge_base_repo

import (
	cache_loader "github.com/XingfenD/yoresee_doc/internal/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type KnowledgeBaseRepository struct {
	db     *gorm.DB
	Loader cache_loader.Loader
}

func NewKnowledgeBaseRepository(db *gorm.DB, redis *redis.Client) *KnowledgeBaseRepository {
	return &KnowledgeBaseRepository{
		db:     db,
		Loader: *cache_loader.NewLoader(redis),
	}
}
