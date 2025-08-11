package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"auction-web/pkg/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProfileController saves the profile into db
func (a *API) ProfileController(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	var request models.Profile
	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind the request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Fetch uuid from token
	tokenUUID := c.GetString("uuid")
	if tokenUUID == "" {
		a.logger.Error("failed to fetch uuid from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not found in token"})
		return
	}

	uid, err := uuid.Parse(tokenUUID)
	if err != nil {
		a.logger.Error("failed to parse the uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid uuid"})
		return
	}
	request.UUID = uid

	err = database.ExecuteQuery(ctx, a.PostgresClient, `INSERT INTO profiles (uuid, first_name, last_name, role, image_url, batting_hand, batting_order, batting_style, bowling_arm, bowling_type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, request.UUID, request.FirstName, request.LastName, request.Role, request.BattingHand, request.BattingOrder, request.BattingStyle, request.BowlingArm, request.BowlingType)
	if err != nil {
		a.logger.Error("failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile updation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": request,
	})
}
