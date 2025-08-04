package middlewares

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(email string) (string, error) {
	// Jwt Secret
	jwtKey := []byte(os.Getenv("TOKEN_SECRET"))

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)

	// Setting data to token
	claims["email"] = email

	// Token expiration time
	claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()

	// Signing the token with the secret key and fetching the token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
