package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (r *postgresRepository) Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	if err := r.db.WithContext(ctx).Create(bookmark).Error; err != nil {
		return nil, dbutils.CatchErr(err)
	}

	return bookmark, nil
}
