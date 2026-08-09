package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/internal/test/data/fixtures"
	jwtMocks "github.com/ManTT-Data/bookmark-management/pkg/jwtutils/mocks"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const otherFixtureUserID = "9ba595b2-345e-4474-ad59-79a7442bc25e"

func TestBookmarkEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := fixtures.NewFixture(t, &fixtures.UserCommonTestCase{})
	require.NoError(t, db.Create(&model.Bookmark{
		ID:          "other-bookmark",
		Code:        "otherbookmark",
		Description: "Other user's bookmark",
		URL:         "https://other.example",
		UserID:      otherFixtureUserID,
	}).Error)

	validatorMock := jwtMocks.NewJWTValidator(t)
	validatorMock.On("ValidateToken", "fixture-token").Return(jwt.MapClaims{"sub": fixtureUserID}, nil)

	engine := api.New(&api.EngineOpts{
		App:    gin.New(),
		Cfg:    &api.Config{},
		Redis:  redisPkg.InitMockRedis(t),
		DB:     db,
		JWTVal: validatorMock,
	})

	created := createBookmark(t, engine, "fixture-token", "Google", "https://google.com")
	bookmarkID := created["id"].(string)
	require.Len(t, created["code"].(string), 10)

	listRequest := authorizedRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=1&sort=created_at", "fixture-token", nil)
	listRecorder := httptest.NewRecorder()
	engine.ServeHTTP(listRecorder, listRequest)
	require.Equal(t, http.StatusOK, listRecorder.Code)

	var listed struct {
		Data       []model.Bookmark `json:"data"`
		Pagination struct {
			Page  int   `json:"page"`
			Limit int   `json:"limit"`
			Total int64 `json:"total"`
		} `json:"pagination"`
	}
	require.NoError(t, json.Unmarshal(listRecorder.Body.Bytes(), &listed))
	require.Len(t, listed.Data, 1)
	require.EqualValues(t, 1, listed.Pagination.Total)
	require.Equal(t, 1, listed.Pagination.Page)
	require.Equal(t, 1, listed.Pagination.Limit)

	updateRequest := authorizedRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, "fixture-token", bytes.NewBufferString(`{"description":"Google Search","url":"https://www.google.com"}`))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRecorder := httptest.NewRecorder()
	engine.ServeHTTP(updateRecorder, updateRequest)
	require.Equal(t, http.StatusOK, updateRecorder.Code)
	require.JSONEq(t, `{"message":"Success"}`, updateRecorder.Body.String())

	foreignUpdateRequest := authorizedRequest(http.MethodPut, "/v1/bookmarks/other-bookmark", "fixture-token", bytes.NewBufferString(`{"description":"No","url":"https://no.example"}`))
	foreignUpdateRequest.Header.Set("Content-Type", "application/json")
	foreignUpdateRecorder := httptest.NewRecorder()
	engine.ServeHTTP(foreignUpdateRecorder, foreignUpdateRequest)
	require.Equal(t, http.StatusNotFound, foreignUpdateRecorder.Code)

	deleteRequest := authorizedRequest(http.MethodDelete, "/v1/bookmarks/"+bookmarkID, "fixture-token", nil)
	deleteRecorder := httptest.NewRecorder()
	engine.ServeHTTP(deleteRecorder, deleteRequest)
	require.Equal(t, http.StatusOK, deleteRecorder.Code)
	require.JSONEq(t, `{"message":"Success"}`, deleteRecorder.Body.String())

	invalidSortRequest := authorizedRequest(http.MethodGet, "/v1/bookmarks?sort=code%3BDROP", "fixture-token", nil)
	invalidSortRecorder := httptest.NewRecorder()
	engine.ServeHTTP(invalidSortRecorder, invalidSortRequest)
	require.Equal(t, http.StatusBadRequest, invalidSortRecorder.Code)

	unauthorizedRequest := httptest.NewRequest(http.MethodGet, "/v1/bookmarks", nil)
	unauthorizedRecorder := httptest.NewRecorder()
	engine.ServeHTTP(unauthorizedRecorder, unauthorizedRequest)
	require.Equal(t, http.StatusUnauthorized, unauthorizedRecorder.Code)
}

func createBookmark(t *testing.T, engine api.Engine, token, description, url string) map[string]any {
	t.Helper()

	request := authorizedRequest(http.MethodPost, "/v1/bookmarks", token, bytes.NewBufferString(`{"description":"`+description+`","url":"`+url+`"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)

	var response struct {
		Message string         `json:"message"`
		Data    map[string]any `json:"data"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
	require.Equal(t, "Create a bookmark successfully!", response.Message)
	return response.Data
}

func authorizedRequest(method, target, token string, body *bytes.Buffer) *http.Request {
	var request *http.Request
	if body == nil {
		request = httptest.NewRequest(method, target, nil)
	} else {
		request = httptest.NewRequest(method, target, body)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	return request
}
