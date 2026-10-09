package cache

import (
	"context"
	"time"
)

// DB stores cached bookmark-list variants grouped by user.
type DB interface {
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error
	DeleteCacheGroup(ctx context.Context, cacheGroupKey string) error
}
