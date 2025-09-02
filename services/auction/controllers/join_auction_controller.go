package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type joinAuctionRequest struct {
	AuctionID string `json:"auction_id"`
}

func (a *API) JoinAuctionController(c *gin.Context) {
	var (
		request  joinAuctionRequest
		response models.Auction
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind join auction request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	email := c.GetString("email")
	if email == "" {
		a.logger.Error("failed to fetch email from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(request.AuctionID)
	if err != nil {
		a.logger.Error("failed to parse auction id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid auction id"})
		return
	}

	filter := bson.M{
		"_id": objectID,
	}
	updateQuery := bson.M{
		"$addToSet": bson.M{
			"joined_by": bson.M{
				"email":         email,
				"is_team_owner": false,
			},
		},
	}

	err = a.MongoDBClient.Collection("auctions").FindOneAndUpdate(ctx, filter, updateQuery).Decode(&response)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Warn("invalid auction id", zap.Any("object_id", request.AuctionID))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Auction not found"})
			return
		}

		a.logger.Error("failed to join auction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}

	// TODO: Delete the cache from redis

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined the auction",
		"auction": response,
	})
}
