package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ManTT-Data/bookmark-management/pkg/jwtutils"
	"github.com/ManTT-Data/bookmark-management/pkg/response"
)

type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

type jwtAuth struct {
	jwtClient jwtutils.JWTValidator
}

func NewJWTAuth(jwtClient jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{jwtClient: jwtClient}
}

func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get token from header
		authHeader := ctx.GetHeader("Authorization")

		//Check if token is empty
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{
				Message: "Unauthorized",
			})
			return
		}

		//Remove Bearer prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")

		//Validate token
		claims, err := j.jwtClient.ValidateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Message{
				Message: "Unauthorized",
			})
			return
		}

		//Set claims to context
		ctx.Set("claims", claims)
		ctx.Next()
	}
}