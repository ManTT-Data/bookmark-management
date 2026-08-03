package requestutils

import (
	"errors"
	"net/http"

	"github.com/ManTT-Data/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidRequest = errors.New("invalid request")
var ErrNoClaims = errors.New("no claims found in request")
var ErrNoUserID = errors.New("no user id found in request")

func GetUserIDFromRequest(ctx *gin.Context) (string, error) {
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{
			Message: "Invalid Token",
		})
		return "", ErrNoClaims
	}
	
	tokenInfo, ok := claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{
			Message: "Invalid Token",
		})
		return "", ErrInvalidRequest
	}

	if userID, ok := tokenInfo["sub"].(string); ok {
		if userID == "" {
			return "", ErrNoUserID
		}
		return userID, nil
	}

	return "", ErrNoUserID
}