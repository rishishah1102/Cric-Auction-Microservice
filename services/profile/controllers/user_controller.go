package controllers

import (
	"github.com/gin-gonic/gin"
)

// UserController fetches the profile and user
func (a *API) UserController(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(c.Request.Context(), constants.DBTimeout)
	// defer cancel()

	// // Fetch uuid from token
	// tokenUUID := c.GetString("uuid")
	// if tokenUUID == "" {
	// 	a.logger.Error("failed to fetch uuid from token")
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "UUID not found in token"})
	// 	return
	// }

	// uid, err := uuid.Parse(tokenUUID)
	// if err != nil {
	// 	a.logger.Error("failed to parse the uuid")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid uuid"})
	// 	return
	// }

}
