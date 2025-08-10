package main

import (
	"auction-web/internal/constants"
	"auction-web/internal/logger"
	"auction-web/internal/router"
	"auction-web/pkg/utils"
	"auction-web/services/auth/controllers"
	"context"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBTimeout)
	defer cancel()

	auctionLogger := logger.Get()
	router := router.NewGinRouter(true)

	api, err := controllers.NewAPI(ctx)
	if err != nil {
		auctionLogger.Error("failed to create API instance", zap.Error(err))
		return
	}
	defer api.PostgresClient.Close()
	defer api.RedisClient.Close()
	api.RegisterRoutes(router)

	utils.StartServer(ctx, router, "auth", "7001")
}
