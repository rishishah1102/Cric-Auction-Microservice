package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) GetAuctionsController(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), constants.DBTimeout)
	// defer cancel()

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the auction micro service",
	})
}
