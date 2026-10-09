package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository/cache"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	bookmarkListCacheGroupKeyFormat = "bookmark:list:%s"
	bookmarkListCacheKeyFormat      = "%d:%d:%s"
	bookmarkListCacheTTL            = 5 * time.Minute
	defaultSort                     = "created_at"
)

type bookmarkServiceCache struct {
	service Service
	cache   cache.DB
}

func NewCache(service Service, cacheDB cache.DB) Service {
	return &bookmarkServiceCache{service: service, cache: cacheDB}
}

func (b *bookmarkServiceCache) List(ctx context.Context, param ListParam) (*ListResult, error) {
	param, err := normalizeCacheListParam(param)
	if err != nil {
		return nil, err
	}

	cacheGroupKey, cacheKey := bookmarkListCacheKey(param)
	data, err := b.cache.GetCacheData(ctx, cacheGroupKey, cacheKey)
	if err == nil {
		var result ListResult
		if err := json.Unmarshal(data, &result); err == nil {
			return &result, nil
		}

		log.Warn().Str("cache_group_key", cacheGroupKey).Msg("Invalid bookmark list cache data")
		b.invalidate(ctx, cacheGroupKey)
	} else if !errors.Is(err, redis.Nil) {
		log.Warn().Err(err).Str("cache_group_key", cacheGroupKey).Msg("Failed to get bookmark list cache")
	}

	result, err := b.service.List(ctx, param)
	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(result)
	if err != nil {
		log.Warn().Err(err).Str("cache_group_key", cacheGroupKey).Msg("Failed to marshal bookmark list cache")
		return result, nil
	}
	if err := b.cache.SetCacheData(ctx, cacheGroupKey, cacheKey, data, bookmarkListCacheTTL); err != nil {
		log.Warn().Err(err).Str("cache_group_key", cacheGroupKey).Msg("Failed to set bookmark list cache")
	}

	return result, nil
}

func (b *bookmarkServiceCache) Delete(ctx context.Context, id, userID string) error {
	if err := b.service.Delete(ctx, id, userID); err != nil {
		return err
	}

	b.invalidateUserBookmarks(ctx, userID)
	return nil
}

func (b *bookmarkServiceCache) Create(ctx context.Context, param CreateParam) (*model.Bookmark, error) {
	result, err := b.service.Create(ctx, param)
	if err != nil {
		return nil, err
	}

	b.invalidateUserBookmarks(ctx, param.UserID)
	return result, nil
}

func (b *bookmarkServiceCache) Update(ctx context.Context, param UpdateParam) error {
	if err := b.service.Update(ctx, param); err != nil {
		return err
	}

	b.invalidateUserBookmarks(ctx, param.UserID)
	return nil
}

func normalizeCacheListParam(param ListParam) (ListParam, error) {
	page, limit, err := normalizePagination(param.Page, param.Limit)
	if err != nil {
		return ListParam{}, err
	}
	if !isSupportedSort(param.Sort) {
		return ListParam{}, ErrInvalidSort
	}

	param.Page = page
	param.Limit = limit
	if param.Sort == "" {
		param.Sort = defaultSort
	}
	return param, nil
}

func bookmarkListCacheKey(param ListParam) (string, string) {
	return fmt.Sprintf(bookmarkListCacheGroupKeyFormat, param.UserID), fmt.Sprintf(
		bookmarkListCacheKeyFormat,
		param.Page,
		param.Limit,
		param.Sort,
	)
}

func (b *bookmarkServiceCache) invalidateUserBookmarks(ctx context.Context, userID string) {
	cacheGroupKey := fmt.Sprintf(bookmarkListCacheGroupKeyFormat, userID)
	b.invalidate(ctx, cacheGroupKey)
}

func (b *bookmarkServiceCache) invalidate(ctx context.Context, cacheGroupKey string) {
	if err := b.cache.DeleteCacheGroup(ctx, cacheGroupKey); err != nil {
		log.Warn().Err(err).Str("cache_group_key", cacheGroupKey).Msg("Failed to invalidate bookmark list cache")
	}
}
