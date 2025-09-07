package controllers

import (
	"auction-web/internal/constants"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type APIrequest struct {
	ID        primitive.ObjectID `json:"team_id"`
	AuctionID primitive.ObjectID `json:"auction_id"`
}

func (a *API) DeleteTeamController(c *gin.Context) {
	var (
		request APIrequest
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind delete team request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	filter := bson.M{
		"_id":        request.ID,
		"auction_id": request.AuctionID,
	}

	_, err := a.MongoDBClient.Collection("teams").DeleteOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Warn("no team found for delete", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found or you are not authorized to delete it"})
			return
		}
		a.logger.Error("failed to delete team in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
		return
	}

	// If team is deleted, we need to delete old data from cache
	cacheKeys := fmt.Sprintf(teamCacheKey, request.AuctionID)
	if _, err = a.RedisClient.Del(ctx, cacheKeys).Result(); err != nil {
		a.logger.Error("failed to delete teams from cache", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Team deleted successfully",
	})
}
