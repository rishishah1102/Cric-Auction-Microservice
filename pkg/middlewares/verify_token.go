package middlewares

import (
	"auction-web/internal/logger"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func VerifyToken(c *gin.Context) {
	auctionLogger := logger.Get()
	jwtKey := []byte(os.Getenv("TOKEN_SECRET"))

	// Fetching token from header of request
	headerToken := c.GetHeader("Authorization")
	if headerToken == "" {
		auctionLogger.Warn("token is required for authentication")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Token is required for authentication",
		})
		return
	}

	// Parse and validate token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(headerToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			auctionLogger.Error("failed to match token sign method", zap.Any("alg", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	// Validating token
	if err != nil {
		auctionLogger.Error("failed to authorize the token", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized, please login again",
		})
		return
	}
	if !token.Valid {
		auctionLogger.Warn("token is invalid")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Store UUID and email in context
	c.Set("uuid", claims.UUID)
	c.Set("email", claims.Email)

	// Token is valid forwarding request
	c.Next()
}
