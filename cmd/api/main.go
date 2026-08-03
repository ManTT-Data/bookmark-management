package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	"github.com/ManTT-Data/bookmark-management/internal/app/model"
	"github.com/ManTT-Data/bookmark-management/pkg/common"
	"github.com/ManTT-Data/bookmark-management/pkg/jwtutils"
	"github.com/ManTT-Data/bookmark-management/pkg/logger"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
	"github.com/ManTT-Data/bookmark-management/pkg/sqldb"
)

// @title           Bookmark API
// @version         1.0
// @description     API documentation for bookmark service.
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Init app config
	cfg := createAPIConfig()

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	logger.SetLogLevel()

	// Init repos
	redisClient := createRedisClient()

	// Init SQL DB
	db, err := createSQLClient()
	common.HandleError(err)

	err = db.AutoMigrate(&model.User{})
	common.HandleError(err)

	// Init JWT Generator & Validator
	jwtGen, err := jwtutils.NewJWTGenerator(cfg.PrivateKeyPath)
	common.HandleError(err)

	jwtVal, err := jwtutils.NewJWTValidator(cfg.PublicKeyPath)
	common.HandleError(err)

	// Init api app
	a := createAPIApp(cfg, redisClient, db, jwtGen, jwtVal)

	// Start api app
	common.HandleError(a.Start())
}

// createRedisClient returns a redis client
func createRedisClient() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	common.HandleError(err)

	return redisClient
}

// createSQLClient returns a gorm DB client
func createSQLClient() (*gorm.DB, error) {
	db, err := sqldb.NewClient("")
	return db, err
}

// createAPIConfig creates API configuration based on environment variables
func createAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	common.HandleError(err)

	return cfg
}

// createAPIApp creates API application based on API configuration
func createAPIApp(cfg *api.Config, redis *redis.Client, db *gorm.DB, jwtGen jwtutils.JWTGenerator, jwtVal jwtutils.JWTValidator) api.Engine {
	app := gin.New()
	a := api.New(&api.EngineOpts{
		App:    app,
		Cfg:    cfg,
		Redis:  redis,
		DB:     db,
		JWTGen: jwtGen,
		JWTVal: jwtVal,
	})

	return a
}
