package user

import (
	"errors"
	"net/http"

	"github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
)

type loginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Login      Login an existing user account
// @Summary      Login an existing user
// @Description  Login an existing user with username and password
// @Tags         user
// @Accept       application/json
// @Produce      application/json
// @Param        input body loginInput true "User login info"
// @Success      200 {object} response.Message "User successfully logged in"
// @Failure      400 {object} response.Message "Invalid input parameters"
// @Failure      500 {object} response.Message "Internal server error"
// @Router       /v1/users/login [post]
func (h *userHandler) Login(ctx *gin.Context) {
	// read input
	input, err := requestutils.BindInputFromRequest[loginInput](ctx)
	if err != nil {
		return
	}

	// call service
	token, err := h.userService.Login(ctx, input.Username, input.Password)
	switch {
	case err == nil:
	case errors.Is(err, user.ErrInvalidCredentials):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Message{
			Message: "Invalid credentials",
		})
		return
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	// return token
	ctx.JSON(http.StatusOK, response.Message{
		Message: token,
	})
}