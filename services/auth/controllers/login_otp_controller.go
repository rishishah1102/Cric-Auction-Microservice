package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/middlewares"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// loginOTPRequest is the struct for login otp controller request body
type loginOTPRequest struct {
	Email string `json:"email"`
	OTP   int    `json:"otp"`
}

// This route is for log in the user and getting a token to make requests
func (a *API) LoginOtpController(c *gin.Context) {
	var request loginOTPRequest

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind login OTP verification request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	val, err := a.RedisClient.Get(ctx, "login_otp:"+request.Email).Result()
	if err == redis.Nil {
		a.logger.Error("failed to fetch the OTP", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired or not found"})
		return
	} else if err != nil {
		a.logger.Error("failed to fetch value from redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	parts := strings.Split(val, ":")
	OTP := parts[0]
	ID := parts[1]

	if strconv.Itoa(request.OTP) != OTP {
		a.logger.Error("failed to validate OTP")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	_ = a.RedisClient.Del(ctx, "login_otp:"+request.Email)

	token, err := middlewares.GenerateToken(ID, request.Email)
	if err != nil {
		a.logger.Error("failed to generate jwt token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User login successful",
		"token":   token,
	})
}
