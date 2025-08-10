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

// API is the struct for all the handlers
type API struct {
	logger         *zap.Logger
	PostgresClient *pgxpool.Pool
	RedisClient    *redis.Client
}

// initTable is the initial query to create the table, partition and index
var initTable = `
-- Create users table
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL NOT NULL PRIMARY KEY,
    email TEXT NOT NULL,
	mobile TEXT NOT NULL,
	image TEXT,
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
) PARTITION BY RANGE (id);

-- Create partition for users table
CREATE TABLE IF NOT EXISTS users_p1 PARTITION OF users FOR VALUES FROM (0) TO (20000);

-- Create index on parent table as it will create indexes on its own for child table
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
`

// NewAPI creates a new API instance
func NewAPI(ctx context.Context) (*API, error) {
	auctionLogger := logger.Get()
	redisCfg := config.LoadRedisConfig()
	postgresCfg := config.LoadPostgresConfig()

	postgresClient, err := database.NewPostgresClient(ctx, postgresCfg.PostgresURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create postgres client")
	}

	if err = database.ExecuteQuery(ctx, postgresClient, initTable); err != nil {
		return nil, logger.WrapError(err, "failed to create table")
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
