package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type playersAPIRequest struct {
	AuctionID primitive.ObjectID `json:"auction_id"`
}

func (a *API) GetAllPlayersController(c *gin.Context) {
	var (
		request playersAPIRequest
		players []models.Player
	)
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind get all players request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	playerType := c.DefaultQuery("role", "all")     // all | Batter | Wicker-Keeper | All-Rounder | Bowler
	playerHammer := c.DefaultQuery("hammer", "all") // all | sold | unsold | upcoming

	filter := bson.M{"auction_id": request.AuctionID}

	if playerType != "all" {
		filter["role"] = playerType
	}

	if playerHammer != "all" {
		filter["hammer"] = playerHammer
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "player_number", Value: 1}}) // Sort Asc

	cursor, err := a.MongoDBClient.Collection("players").Find(ctx, filter, findOptions)
	if err != nil {
		a.logger.Error("failed to fetch players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &players); err != nil {
		a.logger.Error("failed to decode players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while decoding"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Players fetched successfully",
		"auctions": players,
	})
}
