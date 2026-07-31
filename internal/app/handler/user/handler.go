package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ManTT-Data/bookmark-management/internal/app/service/user"
)

type UserHandler interface {
	Register(ctx *gin.Context)
}

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) UserHandler {
	return &userHandler{userService: userService}
}