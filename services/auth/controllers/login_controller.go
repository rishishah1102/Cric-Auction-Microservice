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

// loginRequest is the struct for login controller request body
type loginRequest struct {
	Email string `json:"email"`
}

// This route logs in the user. It takes the email from user and sends otp
func (a *API) LoginController(c *gin.Context) {
	var (
		request loginRequest
		query   = `SELECT id FROM users where email=$1`
	)
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to login bind request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	users, err := database.FetchRecords[models.User](ctx, a.PostgresClient, query, request.Email)
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
	val := strconv.Itoa(otp) + ":" + strconv.FormatInt(users[0].ID, 10)
	if err := a.RedisClient.Set(ctx, "login_otp:"+request.Email, val, TTLTime).Err(); err != nil {
		a.logger.Error("failed to store OTP in redis", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OTP generation failed"})
		return
	}

	go utils.SendEmail(request.Email, "Login OTP", otp)

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent to email",
	})
}
