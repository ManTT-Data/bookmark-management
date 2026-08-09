package bookmark

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

func (s *service) Update(ctx context.Context, param UpdateParam) error {
	return s.repo.Update(ctx, &model.Bookmark{
		ID:          param.ID,
		Description: param.Description,
		URL:         param.URL,
		UserID:      param.UserID,
	})
}
