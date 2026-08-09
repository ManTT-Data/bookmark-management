package bookmark

import (
	"net/http"

	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	bookmarkService "github.com/ManTT-Data/bookmark-management/internal/app/service/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type createInput struct {
	Description string `json:"description" binding:"required,max=255"`
	URL         string `json:"url" binding:"required,url,max=255"`
}

type createResponse struct {
	Message string         `json:"message"`
	Data    model.Bookmark `json:"data"`
}

// Create creates a bookmark for the current user.
// @Summary      Create a bookmark
// @Description  Create a bookmark for the authenticated user
// @Tags         bookmark
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        input body createInput true "Bookmark input"
// @Success      200 {object} createResponse
// @Failure      400 {object} response.Message
// @Failure      401 {object} response.Message
// @Failure      500 {object} response.Message
// @Router       /v1/bookmarks [post]
func (h *handler) Create(ctx *gin.Context) {
	input, err := requestutils.BindInputFromRequest[createInput](ctx)
	if err != nil {
		return
	}

	userID, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	bookmark, err := h.service.Create(ctx, bookmarkService.CreateParam{
		Description: input.Description,
		URL:         input.URL,
		UserID:      userID,
	})
	if err != nil {
		log.Error().Err(err).Str("user_id", userID).Msg("Failed to create bookmark")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, createResponse{
		Message: "Create a bookmark successfully!",
		Data:    *bookmark,
	})
}
