package controllers

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type API struct {
	logger         *zap.Logger
	PostgresClient *pgxpool.Pool
}

// NewAPI creates a new API instance
func NewAPI(ctx context.Context) (*API, error) {
	auctionLogger := logger.Get()
	postgresCfg := config.LoadPostgresConfig()

	postgresClient, err := database.NewPostgresClient(ctx, postgresCfg.PostgresURI)
	if err != nil {
		return nil, logger.WrapError(err, "failed to create postgres client")
	}

	return &API{
		logger:         auctionLogger,
		PostgresClient: postgresClient,
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
