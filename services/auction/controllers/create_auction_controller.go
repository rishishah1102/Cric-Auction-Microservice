package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (a *API) CreateAuctionController(c *gin.Context) {
	var request models.Auction

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind create auction request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	email := c.GetString("email")
	if email == "" {
		a.logger.Error("failed to fetch email from token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	auctionDoc := bson.M{
		"auction_name":   request.AuctionName,
		"auction_image":  request.AuctionImage,
		"auction_date":   request.AuctionDate,
		"created_by":     email,
		"is_ipl_auction": request.IsIPLAuction,
		"joined_by":      []models.JoinedAuctionUser{},
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}

	res, err := a.MongoDBClient.Collection("auctions").InsertOne(ctx, auctionDoc)
	if err != nil {
		a.logger.Error("failed to create auction", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}

	// TODO: Delete the cache from redis

	request.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Auction created successfully",
		"data": map[string]interface{}{
			"auction": request,
		},
	})
}
