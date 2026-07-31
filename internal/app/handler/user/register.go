package user

import (
	"net/http"

	"github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type registerInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register      Register a new user
// @Summary      Register a new user account
// @Description  Create a new user with username, email, and password
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        input body registerInput true "User registration info"
// @Success      200 {object} model.User "User successfully registered"
// @Failure      400 {object} response.Message "Invalid input parameters"
// @Failure      500 {object} response.Message "Internal server error"
// @Router       /v1/users/register [post]
func (h *userHandler) Register(ctx *gin.Context) {
	input := &registerInput{}

	if err := ctx.ShouldBindJSON(input); err != nil {
		ctx.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	param := &user.CreateUserParam{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	createdUser, err := h.userService.CreateUser(ctx, param)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, createdUser)
}