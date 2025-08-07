package utils

import (
	"auction-web/internal/logger"
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartServer(ctx context.Context, router *gin.Engine, service string, port string) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	auctionLogger := logger.Get()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			auctionLogger.Error("failed to start the server", zap.Error(err))
		}
	}()
	auctionLogger.Info("server started")

	<-quit
	auctionLogger.Info("shutting down the server...", zap.String("service", service))

	if err := srv.Shutdown(ctx); err != nil {
		auctionLogger.Error("server forced to shutdwon", zap.Error(err))
	}

	auctionLogger.Info("server exited gracefully")
}
