package url

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage is an interface for URL storage
//
//go:generate mockery --name URLStorage --filename urlstorage.go
type URLStorage interface {
	StoreURL(ctx context.Context, key, value string, exp time.Duration) error
	GetURL(ctx context.Context, key string) (string, error)
}

type urlStorage struct {
	c *redis.Client
}

// NewURLStorage creates a new URL storage
func NewURLStorage(c *redis.Client) URLStorage {
	return &urlStorage{
		c: c,
	}
}

// StoreURL stores a shorten link in cache and return cache key
func (r *urlStorage) StoreURL(ctx context.Context, key, value string, exp time.Duration) error {
	err := r.c.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetURL gets a shorten link from cache
func (r *urlStorage) GetURL(ctx context.Context, key string) (string, error) {
	val, err := r.c.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
