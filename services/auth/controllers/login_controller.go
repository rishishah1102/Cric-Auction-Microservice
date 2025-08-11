package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"auction-web/pkg/models"
	"auction-web/pkg/utils"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// This route logs in the user. It takes the email from user and sends otp
func (a *API) LoginController(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to login bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	users, err := database.FetchRecords[models.User](ctx, a.PostgresClient, `SELECT * FROM users where email=$1`, request.Email)
	if err != nil {
		a.logger.Error("failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	if len(users) == 0 {
		a.logger.Error("failed to find user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	otp := utils.GenerateRandomNumber()
	val := map[string]string{
		"otp":  strconv.Itoa(otp),
		"uuid": users[0].UUID.String(),
	}
	if err := a.RedisClient.HSet(ctx, "login_otp:"+request.Email, val, TTLTime).Err(); err != nil {
		a.logger.Error("failed to store OTP in redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP generation failed"})
		return
	}

	go utils.SendEmail(request.Email, "Login OTP", otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email",
	})
}
