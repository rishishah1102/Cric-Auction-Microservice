package controllers

import (
	"auction-web/internal/database"
	"auction-web/pkg/models"
	"auction-web/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// This route is for getting username and email from frontend and sending otp via email
func (a *API) RegisterController(c *gin.Context) {
	var request models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind register request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if user exists
	users, err := database.FetchRecords[models.User](ctx, a.PostgresClient, `SELECT * FROM users where email=$1`, request.Email)
	if err != nil {
		a.logger.Error("failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	if len(users) != 0 {
		a.logger.Warn("user already exists", zap.String("email", request.Email))
		c.JSON(http.StatusConflict, gin.H{"error": "Account already exists"})
		return
	}

	otp := utils.GenerateRandomNumber()

	if err := a.RedisClient.Set(ctx, "register_otp:"+request.Email, otp, TTLTime).Err(); err != nil {
		a.logger.Error("failed to store OTP in redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP generation failed"})
		return
	}

	go utils.SendEmail(request.Email, "Registration OTP", otp)

	c.JSON(http.StatusCreated, gin.H{
		"message": "OTP sent to email",
	})
}
