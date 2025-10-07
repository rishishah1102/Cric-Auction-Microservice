package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (a *API) DeletePlayerController(c *gin.Context) {
	var request struct {
		PlayerID primitive.ObjectID `json:"player_id" binding:"required"`
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind delete player request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// First get the player to know the auction ID for cache cleanup
	var player models.Player
	err := a.MongoDBClient.Collection("players").FindOne(ctx, bson.M{"_id": request.PlayerID}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
			return
		}
		a.logger.Error("failed to find player for deletion", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find player"})
		return
	}

	// Delete the corresponding match document
	if !player.Match.IsZero() {
		matchFilter := bson.M{"_id": player.Match}
		_, err = a.MongoDBClient.Collection("matches").DeleteOne(ctx, matchFilter)
		if err != nil {
			a.logger.Warn("failed to delete match document for player", zap.Error(err), zap.Any("matchId", player.Match))
		}
	}

	// Delete the player
	_, err = a.MongoDBClient.Collection("players").DeleteOne(ctx, bson.M{"_id": request.PlayerID})
	if err != nil {
		a.logger.Error("failed to delete player", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete player"})
		return
	}

	// Clear relevant cache entries
	cachePattern := fmt.Sprintf("players:auction:%s:*", player.AuctionId.Hex())
	keys, err := a.RedisClient.Keys(ctx, cachePattern).Result()
	if err == nil && len(keys) > 0 {
		if _, err = a.RedisClient.Del(ctx, keys...).Result(); err != nil {
			a.logger.Warn("failed to clear player cache after deletion", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Player deleted successfully",
	})
}
