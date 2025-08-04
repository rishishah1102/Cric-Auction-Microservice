package database

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(uri, password string) (client *redis.Client) {
	client = redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       0,
	})
	return client
}
