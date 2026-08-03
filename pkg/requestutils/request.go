package requestutils

import (
	"errors"
	"io"
	"net/http"

	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	InputValidator = validator.New(validator.WithRequiredStructEnabled())
)

func BindInputFromRequest[T any](ctx *gin.Context) (*T, error) {
	input := new(T)

	// Skip JSON Binding for GET requests to avoid EOF error
	if ctx.Request.Method != http.MethodGet {
		if err := ctx.ShouldBindJSON(input); err != nil && !errors.Is(err, io.EOF) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
			return nil, err
		}
	}

	if err := ctx.ShouldBindUri(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := ctx.ShouldBindQuery(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := ctx.ShouldBindHeader(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	if err := InputValidator.Struct(input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldError(err))
		return nil, err
	}

	return input, nil
}