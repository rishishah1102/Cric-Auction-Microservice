package controllers

import "time"

var (
	TTLTime         = 1 * time.Hour
	auctionCacheKey = "auction_list_%s_%s"
	teamCacheKey    = "team_list_%s"
)
