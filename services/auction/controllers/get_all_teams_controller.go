package controllers

import "github.com/gin-gonic/gin"

func (a *API) GetAllTeamsController(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Welcome to the auction micro service",
	})
}
