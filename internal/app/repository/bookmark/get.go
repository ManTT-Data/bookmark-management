package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
)

func (r *postgresRepository) GetByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&bookmark).Error; err != nil {
		return nil, dbutils.CatchErr(err)
	}

	return &bookmark, nil
}

func (r *postgresRepository) GetByIDAndUserID(ctx context.Context, id, userID string) (*model.Bookmark, error) {
	var bookmark model.Bookmark
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&bookmark).Error; err != nil {
		return nil, dbutils.CatchErr(err)
	}

	return &bookmark, nil
}
