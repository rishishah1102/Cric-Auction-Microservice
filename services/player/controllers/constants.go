package controllers

import "time"

var (
	TTLTime        = 1 * time.Hour
	PlayerCacheKey = "players:auction:%s"
	PlayerTTL      = 5 * time.Minute
)
