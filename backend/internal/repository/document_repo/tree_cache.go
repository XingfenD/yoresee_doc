package document_repo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/key"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	subtreeCacheTTL      = 5 * time.Minute
	subtreeEmptyTTL      = 1 * time.Minute
	subtreeCacheMaxDepth = 3
)

var subtreeCacheSF singleflight.Group

func splitPathPrefixes(path string) []string {
	parts := strings.Split(path, ".")
	out := make([]string, 0, len(parts))
	for i := range parts {
		out = append(out, strings.Join(parts[:i+1], "."))
	}
	return out
}

func getDocPathByID(tx *gorm.DB, docID int64) (string, error) {
	type docPath struct {
		Path string
	}
	var result docPath
	err := tx.Model(&model.Document{}).
		Select("path").
		Where("id = ?", docID).
		Take(&result).Error
	if err != nil {
		return "", err
	}
	return result.Path, nil
}

func (r *DocumentRepository) GetPathByID(id int64) (string, error) {
	return getDocPathByID(storage.DB, id)
}

func (r *DocumentRepository) GetPathByIDWithTx(tx *gorm.DB, id int64) (string, error) {
	return getDocPathByID(tx, id)
}

func getSubtreeVersion(ctx context.Context, path string) (int64, error) {
	if storage.KVS == nil {
		return 0, nil
	}
	val, err := storage.KVS.Get(ctx, key.KeyDocSubtreeVersion(path)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	version, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, nil
	}
	return version, nil
}

func (r *DocumentRepository) BumpSubtreeVersionsByPath(ctx context.Context, path string) error {
	if storage.KVS == nil {
		return nil
	}
	prefixes := splitPathPrefixes(path)
	pipe := storage.KVS.Pipeline()
	for _, prefix := range prefixes {
		pipe.Incr(ctx, key.KeyDocSubtreeVersion(prefix))
	}
	_, err := pipe.Exec(ctx)
	if err == redis.Nil {
		return nil
	}
	return err
}

func getCachedSubtreeIDs(ctx context.Context, key string) ([]int64, bool, error) {
	var ids []int64
	ok, err := cache.GetJSON(ctx, key, &ids)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	return ids, true, nil
}

func setCachedSubtreeIDs(ctx context.Context, key string, ids []int64) {
	ttl := subtreeCacheTTL
	if len(ids) == 0 {
		ttl = subtreeEmptyTTL
	}
	if err := cache.SetJSON(ctx, key, ids, ttl); err != nil {
		logrus.Warnf("set subtree cache failed, key=%s, err=%v", key, err)
	}
}

func fetchDocumentsByIDs(ids []int64) ([]*model.Document, error) {
	if len(ids) == 0 {
		return []*model.Document{}, nil
	}

	var docs []*model.Document
	if err := storage.DB.Model(&model.Document{}).Where("id IN ?", ids).Find(&docs).Error; err != nil {
		return nil, err
	}

	docMap := make(map[int64]*model.Document, len(docs))
	for _, doc := range docs {
		docMap[doc.ID] = doc
	}

	ordered := make([]*model.Document, 0, len(ids))
	for _, id := range ids {
		if doc, ok := docMap[id]; ok {
			ordered = append(ordered, doc)
		}
	}

	return ordered, nil
}
