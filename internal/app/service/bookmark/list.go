package bookmark

import (
	"context"

	bookmarkRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/bookmark"
)

func (s *service) List(ctx context.Context, param ListParam) (*ListResult, error) {
	page, limit, err := normalizePagination(param.Page, param.Limit)
	if err != nil {
		return nil, err
	}
	if !isSupportedSort(param.Sort) {
		return nil, ErrInvalidSort
	}

	bookmarks, total, err := s.repo.ListByUserID(ctx, param.UserID, bookmarkRepo.ListOption{
		Page:  page,
		Limit: limit,
		Sort:  param.Sort,
	})
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Bookmarks: bookmarks,
		Page:      page,
		Limit:     limit,
		Total:     total,
	}, nil
}

func normalizePagination(page, limit int) (int, int, error) {
	if page == 0 {
		page = defaultPage
	}
	if limit == 0 {
		limit = defaultLimit
	}
	if page < 1 || limit < 1 || limit > maxLimit {
		return 0, 0, ErrInvalidPagination
	}

	return page, limit, nil
}

func isSupportedSort(sort string) bool {
	switch sort {
	case "", "created_at", "updated_at", "description", "url":
		return true
	default:
		return false
	}
}
