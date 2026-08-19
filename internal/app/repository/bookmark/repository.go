package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

// ListOption defines filtering, paging, and ordering for user bookmarks.
type ListOption struct {
	Page  int
	Limit int
	Sort  string
}

// Repository defines persistence operations for bookmarks.
//
//go:generate mockery --name Repository --filename repository.go --output ../mocks --outpkg mocks
type Repository interface {
	Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	ListByUserID(ctx context.Context, userID string, option ListOption) ([]model.Bookmark, int64, error)
	GetByCode(ctx context.Context, code string) (*model.Bookmark, error)
	GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bookmark, error)
	Update(ctx context.Context, bookmark *model.Bookmark) error
	Delete(ctx context.Context, id, userID string) error
}
