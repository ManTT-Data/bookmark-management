package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type bookmarkServiceStub struct {
	createFn func(context.Context, CreateParam) (*model.Bookmark, error)
	listFn   func(context.Context, ListParam) (*ListResult, error)
	updateFn func(context.Context, UpdateParam) error
	deleteFn func(context.Context, string, string) error
}

func (s *bookmarkServiceStub) Create(ctx context.Context, param CreateParam) (*model.Bookmark, error) {
	return s.createFn(ctx, param)
}

func (s *bookmarkServiceStub) List(ctx context.Context, param ListParam) (*ListResult, error) {
	return s.listFn(ctx, param)
}

func (s *bookmarkServiceStub) Update(ctx context.Context, param UpdateParam) error {
	return s.updateFn(ctx, param)
}

func (s *bookmarkServiceStub) Delete(ctx context.Context, id, userID string) error {
	return s.deleteFn(ctx, id, userID)
}

type cacheCall struct {
	groupKey string
	cacheKey string
	value    []byte
	exp      time.Duration
}

type cacheStub struct {
	data      map[string][]byte
	getErr    error
	setErr    error
	deleteErr error
	getCalls  []cacheCall
	setCalls  []cacheCall
	deletes   []string
}

func (s *cacheStub) GetCacheData(_ context.Context, groupKey, cacheKey string) ([]byte, error) {
	s.getCalls = append(s.getCalls, cacheCall{groupKey: groupKey, cacheKey: cacheKey})
	if s.getErr != nil {
		return nil, s.getErr
	}
	value, ok := s.data[groupKey+":"+cacheKey]
	if !ok {
		return nil, redis.Nil
	}
	return value, nil
}

func (s *cacheStub) SetCacheData(_ context.Context, groupKey, cacheKey string, value []byte, exp time.Duration) error {
	s.setCalls = append(s.setCalls, cacheCall{groupKey: groupKey, cacheKey: cacheKey, value: value, exp: exp})
	if s.setErr != nil {
		return s.setErr
	}
	if s.data == nil {
		s.data = make(map[string][]byte)
	}
	s.data[groupKey+":"+cacheKey] = value
	return nil
}

func (s *cacheStub) DeleteCacheGroup(_ context.Context, groupKey string) error {
	s.deletes = append(s.deletes, groupKey)
	if s.deleteErr != nil {
		return s.deleteErr
	}
	for key := range s.data {
		if len(key) >= len(groupKey)+1 && key[:len(groupKey)+1] == groupKey+":" {
			delete(s.data, key)
		}
	}
	return nil
}

func TestBookmarkServiceCacheListReturnsCacheHitForCanonicalKey(t *testing.T) {
	ctx := context.Background()
	expected := &ListResult{Bookmarks: []model.Bookmark{{ID: "bookmark-1"}}, Page: 1, Limit: 20, Total: 1}
	data, err := json.Marshal(expected)
	require.NoError(t, err)

	cache := &cacheStub{data: map[string][]byte{"bookmark:list:user-1:1:20:created_at": data}}
	primaryCalls := 0
	service := NewCache(&bookmarkServiceStub{
		listFn: func(context.Context, ListParam) (*ListResult, error) {
			primaryCalls++
			return nil, errors.New("primary service should not be called")
		},
	}, cache)

	result, err := service.List(ctx, ListParam{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, expected, result)
	require.Zero(t, primaryCalls)
	require.Equal(t, []cacheCall{{groupKey: "bookmark:list:user-1", cacheKey: "1:20:created_at"}}, cache.getCalls)
}

func TestBookmarkServiceCacheListCachesMissAndScopesUsers(t *testing.T) {
	ctx := context.Background()
	cache := &cacheStub{data: make(map[string][]byte)}
	primaryCalls := 0
	service := NewCache(&bookmarkServiceStub{
		listFn: func(_ context.Context, param ListParam) (*ListResult, error) {
			primaryCalls++
			return &ListResult{Bookmarks: []model.Bookmark{{Description: param.UserID}}, Page: param.Page, Limit: param.Limit, Total: 1}, nil
		},
	}, cache)

	firstUser, err := service.List(ctx, ListParam{UserID: "user-1"})
	require.NoError(t, err)
	secondUser, err := service.List(ctx, ListParam{UserID: "user-2", Page: 1, Limit: 20, Sort: "created_at"})
	require.NoError(t, err)
	cachedFirstUser, err := service.List(ctx, ListParam{UserID: "user-1", Page: 1, Limit: 20, Sort: "created_at"})
	require.NoError(t, err)

	require.Equal(t, "user-1", firstUser.Bookmarks[0].Description)
	require.Equal(t, "user-2", secondUser.Bookmarks[0].Description)
	require.Equal(t, "user-1", cachedFirstUser.Bookmarks[0].Description)
	require.Equal(t, 2, primaryCalls)
	require.Len(t, cache.setCalls, 2)
	require.Equal(t, bookmarkListCacheTTL, cache.setCalls[0].exp)
}

func TestBookmarkServiceCacheListFallsBackWhenCacheFails(t *testing.T) {
	ctx := context.Background()
	databaseResult := &ListResult{Page: 1, Limit: 20}
	cache := &cacheStub{getErr: errors.New("redis unavailable"), setErr: errors.New("redis unavailable")}
	primaryCalls := 0
	service := NewCache(&bookmarkServiceStub{
		listFn: func(_ context.Context, param ListParam) (*ListResult, error) {
			primaryCalls++
			require.Equal(t, ListParam{UserID: "user-1", Page: 1, Limit: 20, Sort: defaultSort}, param)
			return databaseResult, nil
		},
	}, cache)

	result, err := service.List(ctx, ListParam{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, databaseResult, result)
	require.Equal(t, 1, primaryCalls)
	require.Len(t, cache.setCalls, 1)
}

func TestBookmarkServiceCacheListRemovesCorruptCacheAndFallsBack(t *testing.T) {
	ctx := context.Background()
	cache := &cacheStub{data: map[string][]byte{"bookmark:list:user-1:1:20:created_at": []byte("not-json")}}
	primaryResult := &ListResult{Page: 1, Limit: 20}
	service := NewCache(&bookmarkServiceStub{
		listFn: func(context.Context, ListParam) (*ListResult, error) {
			return primaryResult, nil
		},
	}, cache)

	result, err := service.List(ctx, ListParam{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, primaryResult, result)
	require.Equal(t, []string{"bookmark:list:user-1"}, cache.deletes)
}

func TestBookmarkServiceCacheInvalidatesOnlySuccessfulMutations(t *testing.T) {
	ctx := context.Background()
	testErr := errors.New("mutation failed")

	t.Run("create", func(t *testing.T) {
		cache := &cacheStub{}
		service := NewCache(&bookmarkServiceStub{
			createFn: func(context.Context, CreateParam) (*model.Bookmark, error) { return &model.Bookmark{}, nil },
		}, cache)

		_, err := service.Create(ctx, CreateParam{UserID: "user-1"})
		require.NoError(t, err)
		require.Equal(t, []string{"bookmark:list:user-1"}, cache.deletes)
	})

	t.Run("create error", func(t *testing.T) {
		cache := &cacheStub{}
		service := NewCache(&bookmarkServiceStub{
			createFn: func(context.Context, CreateParam) (*model.Bookmark, error) { return nil, testErr },
		}, cache)

		_, err := service.Create(ctx, CreateParam{UserID: "user-1"})
		require.ErrorIs(t, err, testErr)
		require.Empty(t, cache.deletes)
	})

	t.Run("update", func(t *testing.T) {
		cache := &cacheStub{deleteErr: errors.New("redis unavailable")}
		service := NewCache(&bookmarkServiceStub{
			updateFn: func(context.Context, UpdateParam) error { return nil },
		}, cache)

		err := service.Update(ctx, UpdateParam{UserID: "user-1"})
		require.NoError(t, err)
		require.Equal(t, []string{"bookmark:list:user-1"}, cache.deletes)
	})

	t.Run("update error", func(t *testing.T) {
		cache := &cacheStub{}
		service := NewCache(&bookmarkServiceStub{
			updateFn: func(context.Context, UpdateParam) error { return testErr },
		}, cache)

		err := service.Update(ctx, UpdateParam{UserID: "user-1"})
		require.ErrorIs(t, err, testErr)
		require.Empty(t, cache.deletes)
	})

	t.Run("delete", func(t *testing.T) {
		cache := &cacheStub{}
		service := NewCache(&bookmarkServiceStub{
			deleteFn: func(context.Context, string, string) error { return testErr },
		}, cache)

		err := service.Delete(ctx, "bookmark-1", "user-1")
		require.ErrorIs(t, err, testErr)
		require.Empty(t, cache.deletes)
	})
}
