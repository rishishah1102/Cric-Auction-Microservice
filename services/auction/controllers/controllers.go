package controllers

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// API is the struct for all the handlers
type API struct {
	logger         *zap.Logger
	MongoDB        *mongo.Database
	PostgresClient *pgxpool.Pool
	RedisClient    *redis.Client
}

// NewAPI creates a new API instance
func NewAPI(ctx context.Context) (*API, error) {
	auctionLogger := logger.Get()
	mongoCfg := config.LoadMongoConfig()
	postgresCfg := config.LoadPostgresConfig()
	redisCfg := config.LoadRedisConfig()

	mongoClient, err := database.NewMongoClient(ctx, mongoCfg.MongoURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create mongo client")
	}
	db := mongoClient.Database(mongoCfg.DbName)

	postgresClient, err := database.NewPostgresClient(ctx, postgresCfg.PostgresURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create postgres client")
	}

	redisClient := database.NewRedisClient(redisCfg.RedisURI, redisCfg.RedisPassword)

	return &API{
		logger:         auctionLogger,
		MongoDB:        db,
		PostgresClient: postgresClient,
		RedisClient:    redisClient,
	}, nil
}

// RegisterRoutes register all the handlers to gin router
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Server Check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the auction micro service",
		})
	})
}
