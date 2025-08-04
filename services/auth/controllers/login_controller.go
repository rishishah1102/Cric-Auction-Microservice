package controllers

import (
	"auction-web/pkg/models"
	"auction-web/pkg/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// This route logs in the user. It takes the email from user and sends otp
func (a *API) LoginController(c *gin.Context) {
	var request models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to login bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	err := a.db.Collection("users").FindOne(ctx, bson.M{"email": request.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		a.logger.Error("failed to find user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		a.logger.Error("failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}

	otp := utils.GenerateRandomNumber()
	if err := a.redisClient.Set(ctx, "login_otp:"+request.Email, otp, 5*time.Minute).Err(); err != nil {
		a.logger.Error("failed to store OTP in redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP generation failed"})
		return
	}

	go utils.SendEmail(request.Email, "Login OTP", otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email",
	})
}
