package url

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	urlService "github.com/ManTT-Data/bookmark-management/internal/app/service/url"
	"github.com/ManTT-Data/bookmark-management/pkg/requestutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
)

// ShortenURL represents the Handler for shorten url related operations
type ShortenURL interface {
	ShortenLink(c *gin.Context)
	Redirect(c *gin.Context)
}

type shortenURL struct {
	service urlService.ShortenUrl
}

// NewShortenURL creates a new link Handler
func NewShortenURL(svc urlService.ShortenUrl) ShortenURL {
	return &shortenURL{service: svc}
}

// Redirect 	Forward the request to the original url
// @Tags        link
// @Accept      application/json
// @Produce     application/json
// @Param       code path string true "Shorten code"
// @Success     302
// @Router      /v1/links/redirect/{code} [get]
func (h *shortenURL) Redirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := h.service.GetLinkFromCode(c, code)
	if err != nil {
		if errors.Is(err, urlService.ErrCodeDoesntExist) {
			c.JSON(http.StatusBadRequest, response.InputErrResponse)
			return
		}
		log.Error().Err(err).Str("code", code).Msg("Redirect URL err")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.Redirect(http.StatusFound, url)
}

type shortenResMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type shortenInputBody struct {
	Url string `json:"url" binding:"required,url"`
	Exp int64  `json:"exp" binding:"required,lte=604800"`
}

// ShortenLink   Generate shorten link
// @Summary      Generate shorten url based on original url that last upto 7 days
// @Description  Generate shorten url based on original url that last upto 7 days
// @Tags         link
// @Accept       application/json
// @Produce      application/json
// @Param        input body shortenInputBody true "Input required"
// @Success      200 {object} shortenResMessage
// @Router       /v1/links/shorten [post]
func (h *shortenURL) ShortenLink(c *gin.Context) {
	input, err := requestutils.BindInputFromRequest[shortenInputBody](c)
	if err != nil {
		return
	}

	key, err := h.service.CreateShortenLink(c, input.Url, input.Exp)
	if err != nil {
		log.Error().Err(err).Str("url", input.Url).Int64("exp", input.Exp).Msg("Generate shorten url err")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &shortenResMessage{
		Message: "Shorten URL generated successfully!",
		Code:    key,
	})
}
