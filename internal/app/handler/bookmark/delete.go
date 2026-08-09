package bookmark

import (
	"net/http"

	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
)

// Delete deletes a bookmark owned by the current user.
// @Summary      Delete a bookmark
// @Description  Delete a bookmark owned by the authenticated user
// @Tags         bookmark
// @Security     BearerAuth
// @Produce      application/json
// @Param        id path string true "Bookmark ID"
// @Success      200 {object} response.Message
// @Failure      400 {object} response.Message
// @Failure      401 {object} response.Message
// @Failure      404 {object} response.Message
// @Failure      500 {object} response.Message
// @Router       /v1/bookmarks/{id} [delete]
func (h *handler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	userID, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	if err := h.service.Delete(ctx, id, userID); err != nil {
		handleMutationError(ctx, err, userID, id, "delete")
		return
	}

	ctx.JSON(http.StatusOK, response.Message{Message: "Success"})
}
