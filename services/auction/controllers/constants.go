package controllers

import "time"

var (
	TTLTime  = 1 * time.Hour
	cacheKey = "auction_list_%s_%s"
)
