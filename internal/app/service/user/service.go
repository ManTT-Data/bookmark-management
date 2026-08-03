package user

import (
	"context"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository/user"
	"github.com/ManTT-Data/bookmark-management/pkg/jwtutils"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
)

//go:generate mockery --name Service --filename user.go --output ../mocks --outpkg mocks
// Service defines user business operations.
type Service interface {
	CreateUser(ctx context.Context, param *CreateUserParam) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	Login(ctx context.Context, username string, password string) (string, error)
	GetSelfInfo(ctx context.Context, userId string) (*model.User, error)
}

type userService struct {
	userRepo     user.Repository
	hasher       utils.Hasher
	jwtGenerator jwtutils.JWTGenerator
}

func NewService(userRepo user.Repository, hasher utils.Hasher, jwtGenerator jwtutils.JWTGenerator) Service {
	return &userService{
		userRepo:     userRepo,
		hasher:       hasher,
		jwtGenerator: jwtGenerator,
	}
}
