package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// API is the struct for all the handlers
type API struct {
	logger      *zap.Logger
	db          *mongo.Database
	redisClient *redis.Client
}

// NewAPI creates a new API instance
func NewAPI(logger *zap.Logger, db *mongo.Database, redisClient *redis.Client) *API {
	return &API{
		logger:      logger,
		db:          db,
		redisClient: redisClient,
	}
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
