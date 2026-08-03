package user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

var errServiceFailure = errors.New("service failure")

func newJSONContext(body string) (*gin.Context, *httptest.ResponseRecorder) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return ctx, recorder
}
