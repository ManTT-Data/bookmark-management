package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

func TestRedisDBCacheData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := redisPkg.InitMockRedis(t)
	repo := NewRedisDB(client)
	const (
		groupKey = "bookmark:list:user-1"
		cacheKey = "1:20:created_at"
	)
	value := []byte(`{"bookmarks":[]}`)

	require.NoError(t, repo.SetCacheData(ctx, groupKey, cacheKey, value, 5*time.Minute))

	actual, err := repo.GetCacheData(ctx, groupKey, cacheKey)
	require.NoError(t, err)
	require.Equal(t, value, actual)

	ttl, err := client.TTL(ctx, groupKey).Result()
	require.NoError(t, err)
	require.Positive(t, ttl)
	require.LessOrEqual(t, ttl, 5*time.Minute)

	require.NoError(t, repo.DeleteCacheGroup(ctx, groupKey))
	_, err = repo.GetCacheData(ctx, groupKey, cacheKey)
	require.ErrorIs(t, err, redis.Nil)
}

func TestRedisDBReturnsClosedClientErrors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := redisPkg.InitMockRedis(t)
	require.NoError(t, client.Close())
	repo := NewRedisDB(client)

	require.ErrorIs(t, repo.SetCacheData(ctx, "group", "field", []byte("value"), time.Minute), redis.ErrClosed)
	_, err := repo.GetCacheData(ctx, "group", "field")
	require.ErrorIs(t, err, redis.ErrClosed)
	require.ErrorIs(t, repo.DeleteCacheGroup(ctx, "group"), redis.ErrClosed)
}
