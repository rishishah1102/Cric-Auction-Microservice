package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type joinAuctionRequest struct {
	AuctionID primitive.ObjectID `json:"auction_id"`
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

	// First check if user is already joined
	alreadyJoinedFilter := bson.M{
		"_id":             request.AuctionID,
		"joined_by.email": email,
	}

	count, err := a.MongoDBClient.Collection("auctions").CountDocuments(ctx, alreadyJoinedFilter)
	if err != nil {
		a.logger.Error("failed to check if user already joined", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already joined this auction"})
		return
	}

	filter := bson.M{
		"_id": request.AuctionID,
	}
	updateQuery := bson.M{
		"$addToSet": bson.M{
			"joined_by": bson.M{
				"email": email,
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

	// Clear cache for all auction types for this user
	cacheKeys := []string{
		fmt.Sprintf(auctionCacheKey, "create", email),
		fmt.Sprintf(auctionCacheKey, "join", email),
		fmt.Sprintf(auctionCacheKey, "all", email),
	}
	if _, err = a.RedisClient.Del(ctx, cacheKeys...).Result(); err != nil {
		a.logger.Error("failed to delete join auctions from cache", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully joined the auction",
		"auction": response,
	})
}
