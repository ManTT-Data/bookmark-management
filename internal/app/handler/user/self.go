package user

import (
	"net/http"

	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetSelfInfo  Get self info
// @Summary     Get self info
// @Description Get self info
// @Tags        user
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]interface{}
// @Failure     500 {object} map[string]interface{}
// @Router      /v1/self/info [get]
func (h *userHandler) GetSelfInfo(ctx *gin.Context) {
	userId, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	user, err := h.userService.GetSelfInfo(ctx, userId)
	if err != nil {
		log.Err(err).Msg("Error getting self info")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get self info successfully",
		"data":    user,
	})
}