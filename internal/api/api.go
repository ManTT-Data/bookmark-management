package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ManTT-Data/bookmark-management/docs"
	"github.com/ManTT-Data/bookmark-management/internal/app/handler"
	"github.com/ManTT-Data/bookmark-management/internal/app/repository"
	"github.com/ManTT-Data/bookmark-management/internal/app/service"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
)

// Engine represents the interface for Engine operations
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// api represents the struct for internal API operations
type api struct {
	app   *gin.Engine
	cfg   *Config
	redis *redis.Client
}

// New creates a new API Engine instance
func New(app *gin.Engine, cfg *Config, redis *redis.Client) Engine {
	a := &api{
		app:   app,
		cfg:   cfg,
		redis: redis,
	}
	a.initRoutes()
	return a
}

// Start starts the API Engine
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP to test the API endpoint
func (a *api) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.app.ServeHTTP(w, req)
}

// handlers represents the struct for containing all the necessary handlers for API
type handlers struct {
	healthCheckHandler handler.HealthCheck
	linkHandler        handler.ShortenURL
}

// initHandlers initializes the api handlers
func (a *api) initHandlers() *handlers {
	pingRepo := repository.NewPing(a.redis)
	urlStorageRepo := repository.NewURLStorage(a.redis)
	keyGenerator := utils.NewKeyGenerator()
	shortenUrlService := service.NewShortenUrl(urlStorageRepo, keyGenerator)

	return &handlers{
		healthCheckHandler: handler.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, pingRepo),
		linkHandler:        handler.NewShortenURL(shortenUrlService),
	}
}

// initRoutes initializes the api routes
func (a *api) initRoutes() {
	allHandlers := a.initHandlers()

	// health-check
	a.app.GET("/health-check", allHandlers.healthCheckHandler.HealthCheck)

	// Init swagger routes
	docs.SwaggerInfo.BasePath = a.cfg.BasePath
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Routes := a.app.Group("/v1")
	{
		// link-related
		v1Routes.POST("/links/shorten", allHandlers.linkHandler.ShortenLink)
		v1Routes.GET("/links/redirect/:code", allHandlers.linkHandler.Redirect)
	}
}
