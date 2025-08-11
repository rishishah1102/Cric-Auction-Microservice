package middlewares

import "github.com/golang-jwt/jwt/v4"

// Claims defines the JWT payload structure
type Claims struct {
	Email string `json:"email"`
	UUID  string `json:"uuid"`
	jwt.RegisteredClaims
}
