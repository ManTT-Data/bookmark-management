package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/ManTT-Data/bookmark-management/docs"
	"github.com/ManTT-Data/bookmark-management/internal/api/middleware"
	bookmarkHandler "github.com/ManTT-Data/bookmark-management/internal/app/handler/bookmark"
	healthHandler "github.com/ManTT-Data/bookmark-management/internal/app/handler/healthcheck"
	urlHandler "github.com/ManTT-Data/bookmark-management/internal/app/handler/url"
	userHandler "github.com/ManTT-Data/bookmark-management/internal/app/handler/user"
	bookmarkRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/bookmark"
	cacheRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/cache"
	pingRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/ping"
	urlRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/url"
	userRepo "github.com/ManTT-Data/bookmark-management/internal/app/repository/user"
	bookmarkService "github.com/ManTT-Data/bookmark-management/internal/app/service/bookmark"
	urlService "github.com/ManTT-Data/bookmark-management/internal/app/service/url"
	userService "github.com/ManTT-Data/bookmark-management/internal/app/service/user"
	"github.com/ManTT-Data/bookmark-management/pkg/jwtutils"
	"github.com/ManTT-Data/bookmark-management/pkg/utils"
)

// Engine represents the interface for Engine operations
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// api represents the struct for internal API operations
type api struct {
	app    *gin.Engine
	cfg    *Config
	redis  *redis.Client
	db     *gorm.DB
	jwtGen jwtutils.JWTGenerator
	jwtVal jwtutils.JWTValidator
}

type EngineOpts struct {
	App    *gin.Engine
	Cfg    *Config
	Redis  *redis.Client
	DB     *gorm.DB
	JWTGen jwtutils.JWTGenerator
	JWTVal jwtutils.JWTValidator
}

// New creates a new API Engine instance
func New(opts *EngineOpts) Engine {
	a := &api{
		app:    opts.App,
		cfg:    opts.Cfg,
		redis:  opts.Redis,
		db:     opts.DB,
		jwtGen: opts.JWTGen,
		jwtVal: opts.JWTVal,
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
	bookmarkHandler    bookmarkHandler.Handler
	healthCheckHandler healthHandler.HealthCheck
	linkHandler        urlHandler.ShortenURL
	userHandler        userHandler.UserHandler
}

// initHandlers initializes the api handlers
func (a *api) initHandlers() *handlers {
	pingRepo := pingRepo.NewPing(a.redis)
	urlStorageRepo := urlRepo.NewURLStorage(a.redis)
	keyGenerator := utils.NewKeyGenerator()
	bookmarkRepository := bookmarkRepo.NewRepository(a.db)
	shortenUrlService := urlService.NewShortenUrl(urlStorageRepo, bookmarkRepository, keyGenerator)

	// user
	userRepo := userRepo.NewRepository(a.db)
	hasher := utils.NewHasher()
	userService := userService.NewService(userRepo, hasher, a.jwtGen)

	bookmarkPrimaryService := bookmarkService.NewService(bookmarkRepository, keyGenerator)
	bookmarkCache := cacheRepo.NewRedisDB(a.redis)
	cachedBookmarkService := bookmarkService.NewCache(bookmarkPrimaryService, bookmarkCache)

	return &handlers{
		bookmarkHandler:    bookmarkHandler.NewHandler(cachedBookmarkService),
		healthCheckHandler: healthHandler.NewHealthCheck(a.cfg.ServiceName, a.cfg.InstanceID, pingRepo),
		linkHandler:        urlHandler.NewShortenURL(shortenUrlService),
		userHandler:        userHandler.NewUserHandler(userService),
	}
}

// initRoutes initializes the api routes
func (a *api) initRoutes() {
	a.app.Use(gin.Recovery())
	a.app.Use(middleware.CORS())

	allHandlers := a.initHandlers()

	// health-check
	a.app.GET("/health-check", allHandlers.healthCheckHandler.HealthCheck)

	//Init middleware
	jwtAuth := middleware.NewJWTAuth(a.jwtVal)

	// Init swagger routes
	docs.SwaggerInfo.BasePath = a.cfg.BasePath
	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1Routes := a.app.Group("/v1")
	{
		// link-related
		v1Routes.POST("/links/shorten", allHandlers.linkHandler.ShortenLink)
		v1Routes.GET("/links/redirect/:code", allHandlers.linkHandler.Redirect)

		// user
		v1Routes.POST("/users/register", allHandlers.userHandler.Register)
		v1Routes.POST("/users/login", allHandlers.userHandler.Login)

		// authenticated routes
		v1Routes.Use(jwtAuth.JWTAuth())
		v1Routes.GET("/self/info", allHandlers.userHandler.GetSelfInfo)
		v1Routes.POST("/bookmarks", allHandlers.bookmarkHandler.Create)
		v1Routes.GET("/bookmarks", allHandlers.bookmarkHandler.List)
		v1Routes.PUT("/bookmarks/:id", allHandlers.bookmarkHandler.Update)
		v1Routes.DELETE("/bookmarks/:id", allHandlers.bookmarkHandler.Delete)
	}
}
