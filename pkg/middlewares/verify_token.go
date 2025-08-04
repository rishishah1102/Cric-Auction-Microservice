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

	// Jwt Secret
	jwtKey := []byte(os.Getenv("TOKEN_SECRET"))

	// Fetching token from header of request
	headerToken := c.Request.Header.Get("Authorization")
	if headerToken == "" {
		auctionLogger.Warn("token is required for authentication")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Token is required for authentication",
		})
		return
	}

	// Parse the token
	token, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
		// checking the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			auctionLogger.Error("failed to match token sign method.")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		auctionLogger.Info("token is successfully  parsed")
		return jwtKey, nil
	})

	// Validating token
	if err != nil {
		auctionLogger.Error("failed to authorize the token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized, Please try to login again",
		})
		return
	}
	if !token.Valid {
		auctionLogger.Warn("token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Fetching claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		auctionLogger.Warn("cannot convert token claims to MapClaims")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token is invalid",
		})
		return
	}

	// Extracting email from token claims
	email, ok := claims["email"].(string)
	if !ok {
		auctionLogger.Warn("failed to extract email from token claims")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad token, please login again",
		})
		c.Abort()
		return
	}

	// Store the email in the context for further use
	c.Set("email", email)

	// Token is valid forwarding request
	c.Next()
}
