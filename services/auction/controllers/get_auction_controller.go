package controllers

import (
	"auction-web/internal/constants"
	"auction-web/internal/database"
	"auction-web/pkg/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type auctionAPIRequest struct {
	AuctionID primitive.ObjectID `json:"auction_id"`
}

type auctionUserAPIResp struct {
	models.Auction
	UserNames []userTeam `json:"user_names"`
}

type userTeam struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Pre-compile SQL query to avoid string manipulation overhead
const userQuery = `
	SELECT 
		u.email,
		CONCAT(p.first_name, ' ', p.last_name) AS name
	FROM users u
	INNER JOIN profiles p ON u.id = p.user_id
	WHERE u.email IN (SELECT unnest($1::text[]))
`

func (a *API) GetAuctionController(c *gin.Context) {
	var (
		request  auctionAPIRequest
		response auctionUserAPIResp
	)

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind create team request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var auction models.Auction
	filter := bson.M{"_id": request.AuctionID}
	err := a.MongoDBClient.Collection("auctions").FindOne(ctx, filter).Decode(&auction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Warn("no auction found", zap.Error(err), zap.Any("auction_id", request.AuctionID))
			c.JSON(http.StatusNotFound, gin.H{"error": "Auction not found or you are not authorized to update it"})
			return
		}
		a.logger.Error("failed to find auction in database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find auction"})
		return
	}

	userNames, err := database.FetchRecords[userTeam](ctx, a.PostgresClient, userQuery, auction.JoinedBy)
	if err != nil {
		a.logger.Error("failed to fetch user name from database", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user names from database"})
		return
	}

	response.ID = auction.ID
	response.AuctionName = auction.AuctionName
	response.AuctionImage = auction.AuctionImage
	response.CreatedBy = auction.CreatedBy
	response.AuctionDate = auction.AuctionDate
	response.IsIPLAuction = auction.IsIPLAuction
	response.CreatedAt = auction.CreatedAt
	response.UpdatedAt = auction.UpdatedAt
	response.JoinedBy = append(response.JoinedBy, auction.JoinedBy...)
	response.UserNames = append(response.UserNames, userNames...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Auction fetched successfully",
		"auction": response,
	})
}
