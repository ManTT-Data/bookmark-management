package bookmark

import (
	"net/http"

	bookmarkService "github.com/ManTT-Data/bookmark-management/internal/app/service/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
)

type updateInput struct {
	Description string `json:"description" binding:"required,max=255"`
	URL         string `json:"url" binding:"required,url,max=255"`
}

// Update updates a bookmark owned by the current user.
// @Summary      Update a bookmark
// @Description  Update a bookmark owned by the authenticated user
// @Tags         bookmark
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        id path string true "Bookmark ID"
// @Param        input body updateInput true "Bookmark input"
// @Success      200 {object} response.Message
// @Failure      400 {object} response.Message
// @Failure      401 {object} response.Message
// @Failure      404 {object} response.Message
// @Failure      500 {object} response.Message
// @Router       /v1/bookmarks/{id} [put]
func (h *handler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	input, err := requestutils.BindInputFromRequest[updateInput](ctx)
	if err != nil {
		return
	}

	userID, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	err = h.service.Update(ctx, bookmarkService.UpdateParam{
		ID:          id,
		Description: input.Description,
		URL:         input.URL,
		UserID:      userID,
	})
	if err != nil {
		handleMutationError(ctx, err, userID, id, "update")
		return
	}

	ctx.JSON(http.StatusOK, response.Message{Message: "Success"})
}
