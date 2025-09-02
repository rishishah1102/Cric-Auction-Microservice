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
	"go.uber.org/zap"
)

// registerOTPRequest is the struct for register otp controller request body
type registerOTPRequest struct {
	models.User
	OTP int `json:"otp"`
}

func (a *API) RegisterOtpController(c *gin.Context) {
	var (
		query   = `INSERT INTO users (email, mobile) VALUES ($1, $2) RETURNING id`
		args    []any
		request registerOTPRequest
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind register OTP verification request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	args = []any{request.Email, request.Mobile}

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

	var ID int64
	if err = database.ExecuteQueryReturning(ctx, a.PostgresClient, &ID, query, args...); err != nil {
		a.logger.Error("failed to insert user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}

	if _, err = a.RedisClient.Del(ctx, "register_otp:"+request.Email).Result(); err != nil {
		a.logger.Warn("failed to delete the key from redis", zap.Error(err))
	}

	token, err := middlewares.GenerateToken(strconv.FormatInt(ID, 10), request.Email)
	if err != nil {
		a.logger.Error("failed to generate jwt token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server error from token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": map[string]string{
			"token": token,
		},
	})
}
