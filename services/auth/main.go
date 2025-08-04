package main

import (
	"auction-web/internal/config"
	"auction-web/internal/database"
	"auction-web/internal/logger"
	"auction-web/internal/router"
	"auction-web/services/auth/controllers"
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
)

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg := config.LoadConfig()
	auctionLogger := logger.Get()
	router := router.NewGinRouter(true)

	mongoClient, err := database.NewMongoClient(cfg.MongoDB.MongoURI, cfg.MongoDB.Timeout)
	if err != nil {
		auctionLogger.Error("failed to create mongo client", zap.Error(err))
		return
	}
	defer database.DisconnectMongoClient(mongoClient)
	db := mongoClient.Database(cfg.DbName)

	redisClient := database.NewRedisClient(cfg.Redis.RedisURI, cfg.Redis.RedisPassword)
	defer redisClient.Close()

	api := controllers.NewAPI(auctionLogger, db, redisClient)
	api.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			auctionLogger.Error("failed to start the server", zap.Error(err))
		}
	}()
	auctionLogger.Info("server started")

	<-quit
	auctionLogger.Info("shutting down the server...")

	if err := srv.Shutdown(ctx); err != nil {
		auctionLogger.Error("server forced to shutdwon", zap.Error(err))
	}

	auctionLogger.Info("server exited gracefully")
}
