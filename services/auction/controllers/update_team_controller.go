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

func (a *API) UpdateTeamController(c *gin.Context) {
	var (
		request  models.Team
		response models.Team
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind update team request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	filter := bson.M{
		"_id":        request.ID,
		"auction_id": request.AuctionId,
	}
	update := bson.M{
		"$set": bson.M{
			"team_name":   request.TeamName,
			"team_image":  request.TeamImage,
			"team_owners": request.TeamOwners,
			"updated_at":  time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := a.MongoDBClient.Collection("teams").FindOneAndUpdate(ctx, filter, update, opts).Decode(&response)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Warn("no team found for update", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found or you are not authorized to update it"})
			return
		}
		a.logger.Error("failed to update team in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team"})
		return
	}

	// If team is updated, we need to delete old data from cache
	cacheKeys := fmt.Sprintf(teamCacheKey, request.AuctionId)
	if _, err = a.RedisClient.Del(ctx, cacheKeys).Result(); err != nil {
		a.logger.Error("failed to delete teams from cache", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team updated successfully",
		"team":    response,
	})
}
