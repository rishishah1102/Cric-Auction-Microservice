package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type teamAPIRequest struct {
	AuctionID primitive.ObjectID `json:"auction_id"`
}

func (a *API) GetAllTeamsController(c *gin.Context) {
	var (
		request teamAPIRequest
		teams   []models.Team
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind get all teams request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	teamsKey := fmt.Sprintf(teamCacheKey, request.AuctionID)

	// Try to get from cache first
	val, err := a.RedisClient.Get(ctx, teamsKey).Result()
	if err == nil {
		var cachedTeams []models.Team
		if err = json.Unmarshal([]byte(val), &cachedTeams); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Teams fetched successfully from cache",
				"teams":   cachedTeams,
			})
			return
		} else {
			a.logger.Warn("failed to unmarshal teams from cache", zap.Error(err))
			// Delete invalid cache
			if _, err = a.RedisClient.Del(ctx, teamsKey).Result(); err != nil {
				a.logger.Warn("failed to delete invalid cache key", zap.Error(err))
			}
		}
	}

	filter := bson.M{"auction_id": request.AuctionID}
	cursor, err := a.MongoDBClient.Collection("teams").Find(ctx, filter)
	if err != nil {
		a.logger.Error("failed to fetch teams", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &teams); err != nil {
		a.logger.Error("failed to decode teams", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while decoding"})
		return
	}

	// Cache the results
	if len(teams) > 0 {
		jsonData, err := json.Marshal(teams)
		if err == nil {
			if err = a.RedisClient.Set(ctx, teamsKey, jsonData, TTLTime).Err(); err != nil {
				a.logger.Warn("failed to set teams in redis", zap.Error(err))
			}
		} else {
			a.logger.Warn("failed to marshal teams for caching", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Teams fetched successfully",
		"teams":   teams,
	})
}
