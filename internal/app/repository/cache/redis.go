package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	c *redis.Client
}

func NewRedisDB(c *redis.Client) DB {
	return &redisDB{c: c}
}

func (r *redisDB) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	_, err := r.c.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, cacheGroupKey, cacheKey, value)
		pipe.Expire(ctx, cacheGroupKey, exp)
		return nil
	})
	return err
}

func (r *redisDB) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	return r.c.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
}

func (r *redisDB) DeleteCacheGroup(ctx context.Context, cacheGroupKey string) error {
	return r.c.Del(ctx, cacheGroupKey).Err()
}
