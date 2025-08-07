package controllers

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// API is the struct for all the handlers
type API struct {
	logger      *zap.Logger
	DB          *mongo.Database
	RedisClient *redis.Client
}

// NewAPI creates a new API instance
func NewAPI() (*API, error) {
	auctionLogger := logger.Get()
	mongoCfg := config.LoadMongoConfig()
	redisCfg := config.LoadRedisConfig()

	mongoClient, err := database.NewMongoClient(mongoCfg.MongoURI, mongoCfg.Timeout)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create mongo client")
	}
	db := mongoClient.Database(mongoCfg.DbName)
	auctionLogger.Info("connected to mongo db")

	redisClient := database.NewRedisClient(redisCfg.RedisURI, redisCfg.RedisPassword)

	return &API{
		logger:      auctionLogger,
		DB:          db,
		RedisClient: redisClient,
	}, nil
}

// RegisterRoutes register all the handlers to gin router
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Server Check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the auth micro service",
		})
	})

	authGroup := router.Group("/api/v1/auth")

	// SIGNUP || METHOD POST
	authGroup.POST("/register", a.RegisterController)

	// LOGIN || METHOD POST
	authGroup.POST("/login", a.LoginController)

	// OTP PAGE TO SAVE USER || METHOD POST
	authGroup.POST("/rotp", a.RegisterOtpController)

	// OTP PAGE TO LOGIN USER || METHOD POST
	authGroup.POST("/lotp", a.LoginOtpController)
}
