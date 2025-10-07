package controllers

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// API is the struct for all the handlers
type API struct {
	logger        *zap.Logger
	MongoDBClient *mongo.Database
	RedisClient   *redis.Client
}

// NewAPI creates a new API instance
func NewAPI(ctx context.Context) (*API, error) {
	auctionLogger := logger.Get()
	mongoCfg := config.LoadMongoConfig()
	redisCfg := config.LoadRedisConfig()

	mongoClient, err := database.NewMongoClient(ctx, mongoCfg.MongoURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create mongo client")
	}
	db := mongoClient.Database(mongoCfg.DbName)

	redisClient := database.NewRedisClient(redisCfg.RedisURI, redisCfg.RedisPassword)

	return &API{
		logger:        auctionLogger,
		MongoDBClient: db,
		RedisClient:   redisClient,
	}, nil
}

// RegisterRoutes register all the handlers to gin router
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Server Check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the players micro service",
		})
	})

	playersGroup := router.Group("/api/v1/players")

	playersGroup.POST("/get", a.GetAllPlayersController)

	playersGroup.POST("/save", a.SavePlayerController)

	playersGroup.PATCH("/update", a.UpdatePlayerController)

	playersGroup.DELETE("/delete", a.DeletePlayerController)

	playersGroup.POST("/squad", a.SquadsController)
}
