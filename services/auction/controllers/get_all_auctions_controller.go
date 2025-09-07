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

type auctionAPIResp struct {
	AuctionID    primitive.ObjectID `json:"id"`
	AuctionName  string             `json:"auction_name"`
	AuctionImage string             `json:"auction_image"`
}

func (a *API) GetAllAuctionsController(c *gin.Context) {
	var (
		auctions []models.Auction
		resp     []auctionAPIResp
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	email := c.GetString("email")
	if email == "" {
		a.logger.Error("failed to fetch email from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	auctionType := c.DefaultQuery("type", "all") // all | create | join
	auctionsKey := fmt.Sprintf(auctionCacheKey, auctionType, email)

	// Try to get from cache first
	val, err := a.RedisClient.Get(ctx, auctionsKey).Result()
	if err == nil {
		var cachedAuctions []models.Auction
		if err = json.Unmarshal([]byte(val), &cachedAuctions); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message":  "Auctions fetched successfully from cache",
				"auctions": cachedAuctions,
			})
			return
		} else {
			a.logger.Warn("failed to unmarshal auctions from cache", zap.Error(err))
			// Delete invalid cache
			if _, err = a.RedisClient.Del(ctx, auctionsKey).Result(); err != nil {
				a.logger.Warn("failed to delete invalid cache key", zap.Error(err))
			}
		}
	}

	var filter bson.M
	switch auctionType {
	case "create":
		filter = bson.M{"created_by": email}
	case "join":
		filter = bson.M{"joined_by": email}
	case "all":
		filter = bson.M{
			"$or": []bson.M{
				{"created_by": email},
				{"joined_by": email},
			},
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query type"})
		return
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := a.MongoDBClient.Collection("auctions").Find(ctx, filter, findOptions)
	if err != nil {
		a.logger.Error("failed to fetch auctions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &auctions); err != nil {
		a.logger.Error("failed to decode auctions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while decoding"})
		return
	}

	for _, auction := range auctions {
		var auctionResp auctionAPIResp
		auctionResp.AuctionID = auction.ID
		auctionResp.AuctionName = auction.AuctionName
		auctionResp.AuctionImage = auction.AuctionImage
		resp = append(resp, auctionResp)
	}

	// Cache the results
	if len(resp) > 0 {
		jsonData, err := json.Marshal(resp)
		if err == nil {
			if err = a.RedisClient.Set(ctx, auctionsKey, jsonData, TTLTime).Err(); err != nil {
				a.logger.Warn("failed to set auctions in redis", zap.Error(err))
			}
		} else {
			a.logger.Warn("failed to marshal auctions for caching", zap.Error(err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Auctions fetched successfully",
		"auctions": auctions,
	})
}
