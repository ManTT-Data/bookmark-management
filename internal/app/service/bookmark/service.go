package bookmark

import (
	"context"
	"errors"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	bookmarkRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
)

const (
	bookmarkCodeLength = 8
	maxCodeAttempts    = 5
	defaultPage        = 1
	defaultLimit       = 20
	maxLimit           = 100
)

var (
	ErrInvalidPagination = errors.New("invalid pagination")
	ErrInvalidSort       = errors.New("invalid sort")
	ErrCodeGeneration    = errors.New("could not generate unique bookmark code")
)

// CreateParam contains input required to create a bookmark.
type CreateParam struct {
	Description string
	URL         string
	UserID      string
}

// UpdateParam contains input required to update a bookmark.
type UpdateParam struct {
	ID          string
	Description string
	URL         string
	UserID      string
}

// ListParam contains list filtering and pagination input.
type ListParam struct {
	UserID string
	Page   int
	Limit  int
	Sort   string
}

// ListResult contains the paged bookmarks and matching row count.
type ListResult struct {
	Bookmarks []model.Bookmark
	Page      int
	Limit     int
	Total     int64
}

// Service defines bookmark business operations.
//
//go:generate mockery --name Service --filename service.go --output ../mocks --outpkg mocks
type Service interface {
	Create(ctx context.Context, param CreateParam) (*model.Bookmark, error)
	List(ctx context.Context, param ListParam) (*ListResult, error)
	Update(ctx context.Context, param UpdateParam) error
	Delete(ctx context.Context, id, userID string) error
}

type service struct {
	repo   bookmarkRepo.Repository
	keyGen utils.KeyGenerator
}

func NewService(repo bookmarkRepo.Repository, keyGen utils.KeyGenerator) Service {
	return &service{repo: repo, keyGen: keyGen}
}
