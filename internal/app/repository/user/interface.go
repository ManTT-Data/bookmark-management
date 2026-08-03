package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
)

//go:generate mockery --name Repository --filename user.go --output ../mocks --outpkg mocks
// Repository defines persistence operations for users.
type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}
