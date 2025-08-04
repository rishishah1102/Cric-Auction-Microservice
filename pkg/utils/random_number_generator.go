package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomNumber generates otp
func GenerateRandomNumber() int {
	// Seeding the random number generator with current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generating a random number between 100000 and 999999 (inclusive)
	return rng.Intn(900000) + 100000
}
