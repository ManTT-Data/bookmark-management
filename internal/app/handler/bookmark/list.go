package bookmark

import (
	"errors"
	"net/http"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	bookmarkService "github.com/ManTT-Data/bookmark-management/internal/app/service/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type listInput struct {
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
	Sort  string `form:"sort"`
}

type listResponse struct {
	Data       []model.Bookmark `json:"data"`
	Pagination pagination       `json:"pagination"`
}

type pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// List returns bookmarks for the current user.
// @Summary      List bookmarks
// @Description  List paginated bookmarks for the authenticated user
// @Tags         bookmark
// @Security     BearerAuth
// @Produce      application/json
// @Param        page query int false "Page number, default 1"
// @Param        limit query int false "Items per page, default 20, maximum 100"
// @Param        sort query string false "created_at, updated_at, description, or url"
// @Success      200 {object} listResponse
// @Failure      400 {object} response.Message
// @Failure      401 {object} response.Message
// @Failure      500 {object} response.Message
// @Router       /v1/bookmarks [get]
func (h *handler) List(ctx *gin.Context) {
	input, err := requestutils.BindInputFromRequest[listInput](ctx)
	if err != nil {
		return
	}

	userID, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	result, err := h.service.List(ctx, bookmarkService.ListParam{
		UserID: userID,
		Page:   input.Page,
		Limit:  input.Limit,
		Sort:   input.Sort,
	})
	if err != nil {
		switch {
		case errors.Is(err, bookmarkService.ErrInvalidPagination), errors.Is(err, bookmarkService.ErrInvalidSort):
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		default:
			log.Error().Err(err).Str("user_id", userID).Msg("Failed to list bookmarks")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		}
		return
	}

	ctx.JSON(http.StatusOK, listResponse{
		Data: result.Bookmarks,
		Pagination: pagination{
			Page:  result.Page,
			Limit: result.Limit,
			Total: result.Total,
		},
	})
}
