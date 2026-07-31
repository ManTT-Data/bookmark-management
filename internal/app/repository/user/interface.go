package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}
