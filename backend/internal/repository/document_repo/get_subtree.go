package document_repo

import (
	"context"
	"errors"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"gorm.io/gorm"
)

type DocumentGetSubtreeOperation struct {
	repo         *DocumentRepository
	rootParentID int64
	knowledgeID  *int64
	depth        *int

	directoryOnly bool

	tx *gorm.DB
}

func (r *DocumentRepository) GetSubtree(rootParentID int64) *DocumentGetSubtreeOperation {
	return &DocumentGetSubtreeOperation{
		repo:         r,
		rootParentID: rootParentID,
	}
}

func (op *DocumentGetSubtreeOperation) WithTx(tx *gorm.DB) *DocumentGetSubtreeOperation {
	op.tx = tx
	return op
}

func (op *DocumentGetSubtreeOperation) WithKnowledgeID(knowledgeID *int64) *DocumentGetSubtreeOperation {
	op.knowledgeID = knowledgeID
	return op
}

func (op *DocumentGetSubtreeOperation) WithDepth(depth *int) *DocumentGetSubtreeOperation {
	op.depth = depth
	return op
}

func (op *DocumentGetSubtreeOperation) WithDirectoryOnly(with bool) *DocumentGetSubtreeOperation {
	op.directoryOnly = with
	return op
}

func (op *DocumentGetSubtreeOperation) Exec(ctx context.Context) ([]*model.Document, error) {
	db := storage.DB
	if op.tx != nil {
		db = op.tx
	}

	var documents []*model.Document

	if op.depth != nil && *op.depth == 0 {
		return documents, nil
	}

	type pathDepth struct {
		Path  string
		Depth int
	}
	var root pathDepth
	err := db.Model(&model.Document{}).
		Select("path, depth").
		Where("id = ?", op.rootParentID).
		Take(&root).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return documents, nil
		}
		return nil, err
	}

	// in-transaction query should bypass cache to avoid stale reads
	if op.tx != nil || storage.KVS == nil {
		return op.queryWithRoot(db, root.Path, root.Depth)
	}

	version, err := getSubtreeVersion(ctx, root.Path)
	if err == nil {
		cacheKey := cache.KeyDocSubtree(root.Path, version, op.depth)
		if cachedIDs, ok, err := getCachedSubtreeIDs(ctx, cacheKey); err == nil && ok {
			return fetchDocumentsByIDs(cachedIDs)
		}

		val, err, _ := subtreeCacheSF.Do(cacheKey, func() (interface{}, error) {
			dbDocs, err := op.queryWithRoot(db, root.Path, root.Depth)
			if err != nil {
				return nil, err
			}
			ids := make([]int64, 0, len(dbDocs))
			for _, doc := range dbDocs {
				ids = append(ids, doc.ID)
			}
			setCachedSubtreeIDs(ctx, cacheKey, ids)
			return dbDocs, nil
		})
		if err != nil {
			return nil, err
		}
		if typed, ok := val.([]*model.Document); ok {
			return typed, nil
		}
	}

	return op.queryWithRoot(db, root.Path, root.Depth)
}

func (op *DocumentGetSubtreeOperation) queryWithRoot(db *gorm.DB, rootPath string, rootDepth int) ([]*model.Document, error) {
	var documents []*model.Document

	query := `
		SELECT *
		FROM document_metas
		WHERE deleted_at IS NULL
			AND id <> ?
			AND path <@ ?
	`
	if op.directoryOnly {
		query = `
		SELECT id, external_id, title, parent_id
		FROM document_metas
		WHERE deleted_at IS NULL
			AND id <> ?
			AND path <@ ?
	`
	}
	args := []interface{}{op.rootParentID, rootPath}
	if op.knowledgeID != nil {
		query += " AND knowledge_id = ?"
		args = append(args, *op.knowledgeID)
	}
	if op.depth != nil {
		maxDepth := rootDepth + *op.depth
		query += " AND depth <= ?"
		args = append(args, maxDepth)
	}
	query += " ORDER BY depth, created_at"

	if err := db.Raw(query, args...).Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}
