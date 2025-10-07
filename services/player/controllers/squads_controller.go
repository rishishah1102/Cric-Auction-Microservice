package controllers

import (
	"auction-web/internal/constants"
	"auction-web/pkg/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type SquadResponse struct {
	Batters       []models.Player `json:"batters"`
	Bowlers       []models.Player `json:"bowlers"`
	AllRounders   []models.Player `json:"all_rounders"`
	WicketKeepers []models.Player `json:"wicket_keepers"`
}

func (a *API) SquadsController(c *gin.Context) {
	var request struct {
		PlayerID []primitive.ObjectID `json:"player_id"`
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	defer cancel()

	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("failed to bind squad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	filter := bson.M{"_id": bson.M{"$in": request.PlayerID}}
	cursor, err := a.MongoDBClient.Collection("players").Find(ctx, filter)
	if err != nil {
		a.logger.Error("failed to fetch squad players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error from db"})
		return
	}
	defer cursor.Close(ctx)

	var teamPlayers []models.Player
	if err = cursor.All(ctx, &teamPlayers); err != nil {
		a.logger.Error("failed to decode squad players", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while decoding"})
		return
	}

	response := SquadResponse{
		Batters:       []models.Player{},
		Bowlers:       []models.Player{},
		AllRounders:   []models.Player{},
		WicketKeepers: []models.Player{},
	}
	for _, teamPlayer := range teamPlayers {
		switch teamPlayer.Role {
		case "Batter":
			response.Batters = append(response.Batters, teamPlayer)
		case "Bowler":
			response.Bowlers = append(response.Bowlers, teamPlayer)
		case "All-Rounder":
			response.AllRounders = append(response.AllRounders, teamPlayer)
		case "Wicket-Keeper":
			response.WicketKeepers = append(response.WicketKeepers, teamPlayer)
		default:
			a.logger.Warn("unknown player role",
				zap.String("player_id", teamPlayer.Id.String()),
				zap.String("role", teamPlayer.Role))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Squad fetched successfully",
		"squad":   response,
	})
}
