package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"auction-web/pkg/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProfileController saves the profile into db
func (a *API) ProfileController(c *gin.Context) {
	var (
		query = `WITH upd AS (
			UPDATE profiles
			SET first_name    = $2,
				last_name     = $3,
				role          = $4,
				image_url     = $5,
				batting_hand  = $6,
				batting_order = $7,
				batting_style = $8,
				bowling_arm   = $9,
				bowling_type  = $10,
				updated_at    = now()
			WHERE user_id = $1
			RETURNING *
		)
		INSERT INTO profiles (
			user_id, first_name, last_name, role, image_url, batting_hand, 
			batting_order, batting_style, bowling_arm, bowling_type
		)
		SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		WHERE NOT EXISTS (SELECT 1 FROM upd);
		`
		args []any
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	var request models.Profile
	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind the request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Fetch uuid from token
	tokenID := c.GetString("id")
	if tokenID == "" {
		a.logger.Error("failed to fetch ID from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID not found in token"})
		return
	}

	ID, err := strconv.ParseInt(tokenID, 10, 64)
	if err != nil {
		a.logger.Error("failed to parse the ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ID"})
		return
	}

	request.UserID = ID
	args = []any{
		request.UserID,
		request.FirstName,
		request.LastName,
		request.Role,
		request.ImageURL,
		request.BattingHand,
		request.BattingOrder,
		request.BattingStyle,
		request.BowlingArm,
		request.BowlingType,
	}

	if err = database.ExecuteQuery(ctx, a.PostgresClient, query, args...); err != nil {
		a.logger.Error("failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Profile updation failed"})
		return
	}

	// If profile is updated, we need to delete old data from cache
	if _, err = a.RedisClient.Del(ctx, "auction_profile_"+c.GetString("email")).Result(); err != nil {
		a.logger.Error("failed to delete old data from cache", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": request,
	})
}
