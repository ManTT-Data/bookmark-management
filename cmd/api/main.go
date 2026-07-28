package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/ManTT-Data/bookmark-management/internal/api"
	"github.com/ManTT-Data/bookmark-management/pkg/common"
	"github.com/ManTT-Data/bookmark-management/pkg/logger"
	redisPkg "github.com/ManTT-Data/bookmark-management/pkg/redis"
)

// @title           Bookmark API
// @version         1.0
// @description     API documentation for bookmark service.
// @BasePath        /
func main() {
	// Init app config
	cfg := createAPIConfig()

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	logger.SetLogLevel()

	// Init repos
	redisClient := createRedisClient()

	// Init api app
	a := createAPIApp(cfg, redisClient)

	// Start api app
	common.HandleError(a.Start())
}

// createRedisClient returns a redis client
func createRedisClient() *redis.Client {
	redisClient, err := redisPkg.NewClient("")
	common.HandleError(err)

	return redisClient
}

// createAPIConfig creates API configuration based on environment variables
func createAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	common.HandleError(err)

	return cfg
}

// createAPIApp creates API application based on API configuration
func createAPIApp(cfg *api.Config, redis *redis.Client) api.Engine {
	app := gin.New()
	a := api.New(app, cfg, redis)

	return a
}
