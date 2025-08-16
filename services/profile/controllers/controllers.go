package controllers

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type API struct {
	logger         *zap.Logger
	PostgresClient *pgxpool.Pool
	RedisClient    *redis.Client
}

// NewAPI creates a new API instance
func NewAPI(ctx context.Context) (*API, error) {
	auctionLogger := logger.Get()
	postgresCfg := config.LoadPostgresConfig()
	redisCfg := config.LoadRedisConfig()

	postgresClient, err := database.NewPostgresClient(ctx, postgresCfg.PostgresURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create postgres client")
	}

	redisClient := database.NewRedisClient(redisCfg.RedisURI, redisCfg.RedisPassword)

	return &API{
		logger:         auctionLogger,
		PostgresClient: postgresClient,
		RedisClient:    redisClient,
	}, nil
}

// RegisterRoutes register all the handlers to gin router
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Server Check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the home micro service",
		})
	})

	profileGroup := router.Group("/api/v1")

	profileGroup.POST("/profile", a.ProfileController)

	profileGroup.GET("/profile", a.UserController)
}
