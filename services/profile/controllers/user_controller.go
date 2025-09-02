package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// userProfile struct is the combination of user and profile struct
type userProfile struct {
	Email        string `db:"email" json:"email"`
	Mobile       string `db:"mobile" json:"mobile"`
	FirstName    string `db:"first_name" json:"first_name"`
	LastName     string `db:"last_name" json:"last_name"`
	ImageURL     string `db:"image_url" json:"image_url"`
	Role         string `db:"role" json:"role"`
	BattingHand  string `db:"batting_hand" json:"batting_hand"`
	BattingOrder string `db:"batting_order" json:"batting_order"`
	BattingStyle string `db:"batting_style" json:"batting_style"`
	BowlingArm   string `db:"bowling_arm" json:"bowling_arm"`
	BowlingType  string `db:"bowling_type" json:"bowling_type"`
}

// UserController fetches the profile and user
func (a *API) UserController(c *gin.Context) {
	var query = `
		SELECT 
			u.email,
			u.mobile,
			p.first_name,
			p.last_name,
			p.image_url,
			p.role,
			p.batting_hand,
			p.batting_order,
			p.batting_style,
			p.bowling_arm,
			p.bowling_type
		FROM users u
		INNER JOIN profiles p ON u.id = p.user_id
		WHERE u.email = $1
	`

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	// Fetch uuid from token
	email := c.GetString("email")
	if email == "" {
		a.logger.Error("failed to fetch email from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	userProfileKey := fmt.Sprintf(cacheKey, email)

	val, err := a.RedisClient.Get(ctx, userProfileKey).Result()
	if err == nil {
		var userProfile userProfile
		if err = json.Unmarshal([]byte(val), &userProfile); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message":     "User profile fetched successfully from cache",
				"userProfile": userProfile,
			})
			return
		} else {
			a.logger.Warn("failed to unmarshal user profile", zap.Error(err))
			if _, err = a.RedisClient.Del(ctx, userProfileKey).Result(); err != nil {
				a.logger.Warn("failed to delete the key from redis", zap.Error(err))
			}
		}
	}

	userProfiles, err := database.FetchRecords[userProfile](ctx, a.PostgresClient, query, email)
	if err != nil {
		a.logger.Error("failed to fetch user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}

	if len(userProfiles) == 0 {
		a.logger.Error("failed to find user profile with requested email")
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	userProfileJSON, err := json.Marshal(userProfiles[0])
	if err == nil {
		if err = a.RedisClient.Set(ctx, userProfileKey, userProfileJSON, TTLTime).Err(); err != nil {
			a.logger.Warn("failed to store user profile in redis", zap.Error(err))
		}
	} else {
		a.logger.Warn("failed to marshal user profile", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "User profile fetched successfully",
		"userProfile": userProfiles[0],
	})
}
