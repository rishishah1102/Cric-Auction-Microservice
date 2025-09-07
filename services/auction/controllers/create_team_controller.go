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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (a *API) CreateTeamController(c *gin.Context) {
	var request models.Team

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind create team request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// email := c.GetString("email")
	// if email == "" {
	// 	a.logger.Error("failed to fetch email from token")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
	// 	return
	// }

	teamDoc := bson.M{
		"team_name":   request.TeamName,
		"team_image":  request.TeamImage,
		"auction_id":  request.AuctionId,
		"team_owners": request.TeamOwners,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	res, err := a.MongoDBClient.Collection("teams").InsertOne(ctx, teamDoc)
	if err != nil {
		a.logger.Error("failed to insert team in auction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}

	// Clear cache for all auction types for this user
	cacheKey := fmt.Sprintf(teamCacheKey, request.AuctionId)
	if _, err = a.RedisClient.Del(ctx, cacheKey).Result(); err != nil {
		a.logger.Error("failed to delete teams from cache", zap.Error(err))
	}

	request.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Team inserted successfully",
		"team":    request,
	})
}
