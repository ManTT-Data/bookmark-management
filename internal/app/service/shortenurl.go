package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ManTT-Data/bookmark-management/internal/app/repository"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
)

// ShortenUrl represents the interface of shorten url service
//
//go:generate mockery --name ShortenUrl --filename shortenurl.go
type ShortenUrl interface {
	CreateShortenLink(ctx context.Context, url string, exp int64) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

// service represents the structure of a link service
type shortenURLService struct {
	repo   repository.URLStorage
	keyGen utils.KeyGenerator
}

// NewShortenUrl returns a new healthcheck service
func NewShortenUrl(repo repository.URLStorage, keyGen utils.KeyGenerator) ShortenUrl {
	s := &shortenURLService{
		repo:   repo,
		keyGen: keyGen,
	}

	return s
}

const linkKeyLength = 7

// CreateShortenLink stores a shorten link in cache and return cache key
func (s *shortenURLService) CreateShortenLink(ctx context.Context, url string, exp int64) (string, error) {
	key := s.keyGen.GenerateKey(linkKeyLength)

	res, err := s.repo.GetURL(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	if res != "" {
		return s.CreateShortenLink(ctx, url, exp)
	}

	// put key into redis with value = url
	err = s.repo.StoreURL(ctx, key, url, time.Duration(exp)*time.Second)
	if err != nil {
		return "", err
	}

	return key, nil
}

var ErrCodeDoesntExist = errors.New("code does not exist")

// GetLinkFromCode return the original from shorten code
func (s *shortenURLService) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	link, err := s.repo.GetURL(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeDoesntExist
	}

	return link, err
}
