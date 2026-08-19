package bookmark

import (
	"errors"
	"net/http"

	bookmarkService "github.com/ManTT-Data/bookmark-management/internal/app/service/bookmark"
	"github.com/ManTT-Data/bookmark-management/pkg/dbutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Handler handles authenticated bookmark requests.
type Handler interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type handler struct {
	service bookmarkService.Service
}

func NewHandler(service bookmarkService.Service) Handler {
	return &handler{service: service}
}

func handleMutationError(ctx *gin.Context, err error, userID, id, operation string) {
	if errors.Is(err, dbutils.ErrRecordNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.Message{Message: "Bookmark not found"})
		return
	}

	log.Error().Err(err).Str("user_id", userID).Str("bookmark_id", id).Msg("Failed to " + operation + " bookmark")
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalErrResponse)
}
