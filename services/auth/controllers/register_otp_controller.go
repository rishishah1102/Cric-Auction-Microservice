package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"auction-web/pkg/middlewares"
	"auction-web/pkg/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (a *API) RegisterOtpController(c *gin.Context) {
	var request struct {
		models.User
		OTP int `json:"otp"`
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind register OTP verification request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	storedOtp, err := a.RedisClient.Get(ctx, "register_otp:"+request.Email).Result()
	if err == redis.Nil {
		a.logger.Error("failed to fetch the OTP", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
		return
	} else if err != nil {
		a.logger.Error("failed to fetch value from redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from caching"})
		return
	}

	if strconv.Itoa(request.OTP) != storedOtp {
		a.logger.Error("failed to validate OTP")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	var uuid uuid.UUID
	if err = database.ExecuteQueryReturning(ctx, a.PostgresClient, &uuid, `INSERT INTO users (email, mobile) VALUES ($1, $2, $3) RETURNING uuid`, request.Email, request.Mobile); err != nil {
		a.logger.Error("failed to insert user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	_ = a.RedisClient.Del(ctx, "register_otp:"+request.Email)

	token, err := middlewares.GenerateToken(uuid.String(), request.Email)
	if err != nil {
		a.logger.Error("failed to generate jwt token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server error from token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}
