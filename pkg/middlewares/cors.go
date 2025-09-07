package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {
	// Allowing multiple origins
	origin := c.Request.Header.Get("Origin")
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
		"http://192.168.29.239:3000",
		// TODO
	}

	for _, o := range allowedOrigins {
		if origin == o {
			c.Header("Access-Control-Allow-Origin", origin) // CORS
			break
		}
	}

	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
	c.Header("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
