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

func (a *API) UpdateAuctionController(c *gin.Context) {
	var (
		request  models.Auction
		response models.Auction
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind update auction request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	email := c.GetString("email")
	if email == "" {
		a.logger.Error("failed to fetch email from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	filter := bson.M{
		"_id":        request.ID,
		"created_by": email,
	}
	update := bson.M{
		"$set": bson.M{
			"auction_name":   request.AuctionName,
			"auction_image":  request.AuctionImage,
			"auction_date":   request.AuctionDate,
			"is_ipl_auction": request.IsIPLAuction,
			"updated_at":     time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := a.MongoDBClient.Collection("auctions").FindOneAndUpdate(ctx, filter, update, opts).Decode(&response)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Warn("no auction found for update", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found or you are not authorized to update it"})
			return
		}
		a.logger.Error("failed to update auction in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update auction"})
		return
	}

	// If auction is updated, we need to delete old data from cache
	cacheKeys := []string{
		fmt.Sprintf(auctionCacheKey, "create", email),
		fmt.Sprintf(auctionCacheKey, "join", email),
		fmt.Sprintf(auctionCacheKey, "all", email),
	}
	if _, err = a.RedisClient.Del(ctx, cacheKeys...).Result(); err != nil {
		a.logger.Error("failed to delete auctions from cache", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Auction updated successfully",
		"auction": response,
	})
}
