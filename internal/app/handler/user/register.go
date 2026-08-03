package user

import (
	"errors"
	"net/http"

	"github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
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
	input, err := requestutils.BindInputFromRequest[registerInput](ctx)
	if err != nil {
		return
	}

	param := &user.CreateUserParam{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	createdUser, err := h.userService.CreateUser(ctx, param)

	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, createdUser)
	case errors.Is(err, dbutils.ErrDuplicateEmail):
		ctx.JSON(http.StatusBadRequest, response.Message{
			Message: "Email already exists",
		})
	case errors.Is(err, dbutils.ErrDuplicateUserName):
		ctx.JSON(http.StatusBadRequest, response.Message{
			Message: "Username already exists",
		})
	default:
		log.Error().Err(err).Msg("Failed to create user")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
	}
}