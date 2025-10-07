package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (a *API) UpdatePlayerController(c *gin.Context) {
	var player models.Player

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&player); err != nil {
		a.logger.Error("failed to bind update player request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if player.Id.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player ID is required"})
		return
	}

	// Set updated timestamp
	player.UpdatedAt = time.Now()

	filter := bson.M{"_id": player.Id}
	replaceOptions := options.FindOneAndReplace().SetReturnDocument(options.After)

	var updatedPlayer models.Player
	err := a.MongoDBClient.Collection("players").FindOneAndReplace(ctx, filter, player, replaceOptions).Decode(&updatedPlayer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		a.logger.Error("failed to update player", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
		return
	}

	// Clear relevant cache entries
	cachePattern := fmt.Sprintf("players:auction:%s", player.AuctionId.Hex())
	keys, err := a.RedisClient.Keys(ctx, cachePattern).Result()
	if err == nil && len(keys) > 0 {
		if _, err = a.RedisClient.Del(ctx, keys...).Result(); err != nil {
			a.logger.Warn("failed to clear player cache after update", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Player updated successfully",
		"player":  updatedPlayer,
	})
}
