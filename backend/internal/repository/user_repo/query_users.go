package user_repo

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/XingfenD/yoresee_doc/internal/model"
	"github.com/XingfenD/yoresee_doc/pkg/cache"
	"github.com/XingfenD/yoresee_doc/pkg/storage"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QueryUsersOperation struct {
	repo     *UserRepository
	tx       *gorm.DB
	keyword  *string
	userIDs  []int64
	page     int
	pageSize int
}

type userQueryCache struct {
	UserIDs []int64 `json:"user_ids"`
	Total   int64   `json:"total"`
}

const userQueryCacheTTL = time.Minute

func (r *UserRepository) QueryUsers() *QueryUsersOperation {
	return &QueryUsersOperation{
		repo: r,
	}
}

func (op *QueryUsersOperation) WithTx(tx *gorm.DB) *QueryUsersOperation {
	op.tx = tx
	return op
}

func (op *QueryUsersOperation) WithKeyword(keyword *string) *QueryUsersOperation {
	op.keyword = keyword
	return op
}

func (op *QueryUsersOperation) WithUserIDs(userIDs []int64) *QueryUsersOperation {
	op.userIDs = userIDs
	return op
}

func (op *QueryUsersOperation) WithPagination(page, pageSize int) *QueryUsersOperation {
	op.page = page
	op.pageSize = pageSize
	return op
}

func (op *QueryUsersOperation) ExecWithTotal() ([]model.User, int64, error) {
	if op.tx == nil {
		op.tx = storage.DB
	}

	if op.tx == storage.DB && op.page > 0 && op.pageSize > 0 {
		if users, total, ok := op.tryCache(context.Background()); ok {
			return users, total, nil
		}
	}

	query := op.tx.Model(&model.User{})
	if len(op.userIDs) > 0 {
		query = query.Where("id IN ?", op.userIDs)
	}
	if op.keyword != nil {
		trimmed := strings.TrimSpace(*op.keyword)
		if trimmed != "" {
			like := "%" + trimmed + "%"
			query = query.Where("username ILIKE ? OR email ILIKE ? OR external_id ILIKE ?", like, like, like)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if op.page > 0 && op.pageSize > 0 {
		offset := (op.page - 1) * op.pageSize
		query = query.Offset(offset).Limit(op.pageSize)
	}

	var users []model.User
	if err := query.Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	if op.tx == storage.DB && op.page > 0 && op.pageSize > 0 {
		op.storeCache(context.Background(), users, total)
	}

	return users, total, nil
}

func (op *QueryUsersOperation) tryCache(ctx context.Context) ([]model.User, int64, bool) {
	version := getUserQueryVersion(ctx)
	queryHash := buildUserQueryHash(op.keyword, op.userIDs)
	cacheKey := cache.KeyUserQueryList(fmt.Sprintf("%d:%s", version, queryHash), op.page, op.pageSize)

	var cached userQueryCache
	ok, err := cache.GetJSON(ctx, cacheKey, &cached)
	if err != nil || !ok {
		if err != nil {
			logrus.Warnf("user query cache read failed: %v", err)
		}
		return nil, 0, false
	}

	if len(cached.UserIDs) == 0 {
		return []model.User{}, cached.Total, true
	}

	userMap, err := op.repo.MGetUserByID(cached.UserIDs).WithTx(op.tx).Exec()
	if err != nil {
		logrus.Warnf("user query cache mget failed: %v", err)
		return nil, 0, false
	}

	users := make([]model.User, 0, len(cached.UserIDs))
	for _, id := range cached.UserIDs {
		if user, ok := userMap[id]; ok {
			users = append(users, *user)
		}
	}

	return users, cached.Total, true
}

func (op *QueryUsersOperation) storeCache(ctx context.Context, users []model.User, total int64) {
	userIDs := make([]int64, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	version := getUserQueryVersion(ctx)
	queryHash := buildUserQueryHash(op.keyword, op.userIDs)
	cacheKey := cache.KeyUserQueryList(fmt.Sprintf("%d:%s", version, queryHash), op.page, op.pageSize)
	if err := cache.SetJSON(ctx, cacheKey, &userQueryCache{UserIDs: userIDs, Total: total}, userQueryCacheTTL); err != nil {
		logrus.Warnf("user query cache set failed: %v", err)
	}
}

func buildUserQueryHash(keyword *string, userIDs []int64) string {
	var trimmed string
	if keyword != nil {
		trimmed = strings.TrimSpace(*keyword)
	}

	sortedIDs := append([]int64{}, userIDs...)
	sort.Slice(sortedIDs, func(i, j int) bool { return sortedIDs[i] < sortedIDs[j] })

	builder := strings.Builder{}
	builder.WriteString(trimmed)
	builder.WriteString("|")
	for i, id := range sortedIDs {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%d", id))
	}

	sum := sha1.Sum([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}

func getUserQueryVersion(ctx context.Context) int64 {
	val, err := storage.KVS.Get(ctx, cache.KeyUserQueryVersion()).Int64()
	if err != nil {
		return 1
	}
	if val <= 0 {
		return 1
	}
	return val
}
