package middlewares

import "github.com/golang-jwt/jwt/v4"

// Claims defines the JWT payload structure
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}
