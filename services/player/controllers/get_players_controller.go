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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// GetAllPlayersController with Redis caching
func (a *API) GetAllPlayersController(c *gin.Context) {
	var request struct {
		AuctionID primitive.ObjectID `json:"auction_id" binding:"required"`
	}
	var players []models.Player

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind get all players request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Generate cache key
	cacheKey := fmt.Sprintf(PlayerCacheKey, request.AuctionID.Hex())

	val, err := a.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedPlayers []models.Player
		if err = json.Unmarshal([]byte(val), &cachedPlayers); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Players fetched successfully from cache",
				"players": cachedPlayers,
			})
			return
		} else {
			a.logger.Warn("failed to unmarshal players from cache", zap.Error(err))
			// Delete invalid cache
			if _, err = a.RedisClient.Del(ctx, cacheKey).Result(); err != nil {
				a.logger.Warn("failed to delete invalid cache key", zap.Error(err))
			}
		}
	}

	filter := bson.M{"auction_id": request.AuctionID}

	findOptions := options.Find().SetSort(bson.D{{Key: "player_number", Value: 1}}) // Sort Asc

	cursor, err := a.MongoDBClient.Collection("players").Find(ctx, filter, findOptions)
	if err != nil {
		a.logger.Error("failed to fetch players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &players); err != nil {
		a.logger.Error("failed to decode players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while decoding"})
		return
	}

	// Cache the results
	if len(players) > 0 {
		jsonData, err := json.Marshal(players)
		if err == nil {
			if err = a.RedisClient.Set(ctx, cacheKey, jsonData, PlayerTTL).Err(); err != nil {
				a.logger.Warn("failed to set players in redis", zap.Error(err))
			}
		} else {
			a.logger.Warn("failed to marshal players for caching", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Players fetched successfully",
		"players": players,
	})
}
